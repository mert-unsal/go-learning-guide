#!/usr/bin/env python3
"""
md2pdf — Professional PDF Generator from Markdown Files

Converts a collection of Markdown files into a single, professionally formatted PDF
with all the features expected of a modern digital book:

  • Clickable Table of Contents with page numbers and dot leaders
  • PDF Bookmarks (sidebar outline) with Part → Chapter hierarchy
  • Running headers (chapter name) on every content page
  • Page numbers on every page (except cover)
  • "How to Use This Book" guide page (auto-generated or custom)
  • Professional typography (12pt body, 1.5 line-height, ~70 char lines)
  • break-inside: avoid — code blocks, tables, and diagrams stay together
  • Proper PDF metadata (title, author, subject, keywords)
  • Font embedding (via Chrome) and optimized file size

Pipeline:
  1. Read Markdown files + optional config (JSON)
  2. Build full HTML document (cover, how-to-use, TOC, chapters)
  3. Chrome Headless pass 1 → temp PDF → extract real page numbers
  4. Chrome Headless pass 2 → PDF with accurate TOC page numbers
  5. PyMuPDF post-process → bookmarks, running headers, footers, metadata

Requirements:
  pip install markdown pymupdf
  Google Chrome or Microsoft Edge (for headless PDF rendering)

Usage:
  # With a config file (recommended for structured books):
  python tools/md2pdf.py --config tools/book_config.json

  # Quick mode (auto-detect chapters from numbered .md files):
  python tools/md2pdf.py --source-dir learnings --title "My Book" --output book.pdf

  # Show all options:
  python tools/md2pdf.py --help

Config file format (JSON):
  See tools/book_config.json for a complete example.

Lessons learned (from iterating on this tool):
  - Chrome's --no-pdf-header-footer AND --print-to-pdf-no-header both needed
  - Set <title> to empty to prevent Chrome from using it as a header
  - PyMuPDF's built-in Helvetica ("helv") can't render em-dashes — use ASCII dashes
  - On Windows, save PyMuPDF output to a NEW file (don't overwrite the open file)
  - Two-pass rendering is required for accurate TOC page numbers
  - CSS break-inside:avoid on <pre> keeps code blocks together, but allow breaks
    for blocks >55 lines (they can't fit on one page anyway)
  - Chrome's text extraction via PyMuPDF shows overlaid text at end of content
    stream, not in visual position order — this is normal
  - Newlines in PDF text extraction break substring search — normalize first
"""

import argparse
import json
import os
import re
import subprocess
import sys

try:
    import markdown
    from markdown.extensions.tables import TableExtension
    from markdown.extensions.fenced_code import FencedCodeExtension
except ImportError:
    print("ERROR: 'markdown' package not found. Install with: pip install markdown")
    sys.exit(1)

try:
    import fitz  # PyMuPDF
except ImportError:
    print("ERROR: 'pymupdf' package not found. Install with: pip install pymupdf")
    sys.exit(1)


# ═══════════════════════════════════════════════════════════════════════
# Constants
# ═══════════════════════════════════════════════════════════════════════

PAGE_W = 595.28  # A4 width in points
PAGE_H = 841.89  # A4 height in points
MM = 2.835       # 1mm in points

# Chrome paths to try (Windows + macOS + Linux)
CHROME_PATHS = [
    r"C:\Program Files\Google\Chrome\Application\chrome.exe",
    r"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe",
    r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
    r"C:\Program Files\Microsoft\Edge\Application\msedge.exe",
    "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
    "/usr/bin/google-chrome",
    "/usr/bin/chromium-browser",
    "/usr/bin/chromium",
]


# ═══════════════════════════════════════════════════════════════════════
# Configuration
# ═══════════════════════════════════════════════════════════════════════

class BookConfig:
    """Book configuration — loaded from JSON or built from CLI args."""

    def __init__(self):
        self.title = "Untitled Book"
        self.subtitle = ""
        self.author = ""
        self.description = ""
        self.subject = ""
        self.keywords = ""
        self.source_dir = "."
        self.output = "book.pdf"
        self.howto_text = ""       # Custom "How to Use" content (markdown)
        self.parts = []            # [{"name": "Part I", "chapters": [{"file", "short"}, ...]}]
        self.chrome_path = ""

        # Layout
        self.margin_top_mm = 32
        self.margin_bottom_mm = 25
        self.margin_lr_mm = 30
        self.body_font_size_pt = 12
        self.code_font_size_pt = 9
        self.line_height = 1.5
        self.header_font_size_pt = 8.5
        self.footer_font_size_pt = 9

        # Header/footer placement (mm from top of page)
        self.header_text_y_mm = 17
        self.header_line_y_mm = 23
        self.footer_text_y_mm = 291

    @classmethod
    def from_json(cls, path):
        """Load config from a JSON file."""
        cfg = cls()
        with open(path, "r", encoding="utf-8") as f:
            data = json.load(f)

        for key in ["title", "subtitle", "author", "description", "subject",
                     "keywords", "source_dir", "output", "howto_text", "chrome_path"]:
            if key in data:
                setattr(cfg, key, data[key])

        layout = data.get("layout", {})
        for key in ["margin_top_mm", "margin_bottom_mm", "margin_lr_mm",
                     "body_font_size_pt", "code_font_size_pt", "line_height",
                     "header_font_size_pt", "footer_font_size_pt",
                     "header_text_y_mm", "header_line_y_mm", "footer_text_y_mm"]:
            if key in layout:
                setattr(cfg, key, layout[key])

        if "parts" in data:
            cfg.parts = data["parts"]

        # Resolve source_dir relative to config file
        config_dir = os.path.dirname(os.path.abspath(path))
        if not os.path.isabs(cfg.source_dir):
            cfg.source_dir = os.path.normpath(os.path.join(config_dir, cfg.source_dir))
        if not os.path.isabs(cfg.output):
            cfg.output = os.path.normpath(os.path.join(config_dir, cfg.output))

        return cfg

    @classmethod
    def from_args(cls, args):
        """Build config from CLI arguments."""
        cfg = cls()
        cfg.title = args.title or "Untitled Book"
        cfg.subtitle = args.subtitle or ""
        cfg.author = args.author or ""
        cfg.source_dir = os.path.abspath(args.source_dir)
        cfg.output = os.path.abspath(args.output)
        cfg.chrome_path = args.chrome or ""
        return cfg

    def auto_detect_chapters(self):
        """Auto-detect chapters from numbered .md files if no parts defined."""
        if self.parts:
            return

        md_files = sorted([
            f for f in os.listdir(self.source_dir)
            if f.endswith(".md") and f.lower() != "readme.md"
        ])

        if not md_files:
            print(f"ERROR: No .md files found in {self.source_dir}")
            sys.exit(1)

        chapters = []
        for fname in md_files:
            # Extract number prefix and short title from filename
            match = re.match(r"^(\d+)[_\-](.+)\.md$", fname)
            if match:
                num = match.group(1)
                raw_title = match.group(2)
            else:
                num = str(len(chapters) + 1).zfill(2)
                raw_title = fname.replace(".md", "")

            # Convert filename to readable title
            short = raw_title.replace("_", " ").replace("-", " ").title()

            chapters.append({
                "file": fname,
                "short": short,
                "num": num,
                "id": f"ch{num}",
            })

        # Put all chapters in a single part
        self.parts = [{"name": "Chapters", "chapters": chapters}]

    def resolve_chapters(self):
        """Ensure all chapters have 'id' and 'num' fields."""
        counter = 1
        for part in self.parts:
            for ch in part["chapters"]:
                if "num" not in ch:
                    ch["num"] = str(counter).zfill(2)
                if "id" not in ch:
                    ch["id"] = f"ch{ch['num']}"
                counter += 1


# ═══════════════════════════════════════════════════════════════════════
# CSS Stylesheet (parameterized)
# ═══════════════════════════════════════════════════════════════════════

def build_css(cfg: BookConfig) -> str:
    return f"""
@page {{
    size: A4;
    margin: {cfg.margin_top_mm}mm {cfg.margin_lr_mm}mm {cfg.margin_bottom_mm}mm {cfg.margin_lr_mm}mm;
}}

* {{ box-sizing: border-box; }}

body {{
    font-family: 'Segoe UI', 'Helvetica Neue', Arial, sans-serif;
    font-size: {cfg.body_font_size_pt}pt;
    line-height: {cfg.line_height};
    color: #1a1a1a;
    margin: 0;
    padding: 0;
}}

/* ── Cover ─────────────────────────────────────────────────── */
.cover {{
    page-break-after: always;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    min-height: 80vh;
    text-align: center;
}}
.cover h1 {{
    font-size: 36pt;
    color: #1a56db;
    margin-bottom: 6px;
    border: none;
    padding: 0;
    letter-spacing: -0.5px;
}}
.cover .subtitle {{
    font-size: 15pt;
    color: #555;
    margin-bottom: 36px;
    font-style: italic;
}}
.cover .tagline {{
    font-size: 11pt;
    color: #444;
    max-width: 460px;
    line-height: 1.9;
}}

/* ── Page Breaks ───────────────────────────────────────────── */
.chapter-break {{ page-break-before: always; }}

/* ── Headings ──────────────────────────────────────────────── */
h1 {{
    font-size: 21pt;
    color: #1a56db;
    border-bottom: 3px solid #1a56db;
    padding-bottom: 6px;
    margin-top: 30px;
    margin-bottom: 16px;
    break-after: avoid;
    line-height: 1.3;
}}
h2 {{
    font-size: 16pt;
    color: #2563eb;
    border-bottom: 1px solid #ddd;
    padding-bottom: 4px;
    margin-top: 24px;
    margin-bottom: 12px;
    break-after: avoid;
    line-height: 1.3;
}}
h3 {{
    font-size: 13pt;
    color: #3b82f6;
    margin-top: 20px;
    margin-bottom: 8px;
    break-after: avoid;
}}
h4 {{
    font-size: 11.5pt;
    color: #60a5fa;
    margin-top: 16px;
    margin-bottom: 6px;
    break-after: avoid;
}}

/* ── Code Blocks — KEEP TOGETHER ───────────────────────────── */
pre {{
    background: #f6f8fa;
    border: 1px solid #d0d7de;
    border-left: 3px solid #3b82f6;
    border-radius: 6px;
    padding: 12px 14px;
    font-family: 'Cascadia Code', 'Consolas', 'Courier New', monospace;
    font-size: {cfg.code_font_size_pt}pt;
    line-height: 1.45;
    white-space: pre-wrap;
    word-wrap: break-word;
    break-inside: avoid;
    page-break-inside: avoid;
    margin: 10px 0;
}}
pre.allow-break {{
    break-inside: auto;
    page-break-inside: auto;
}}
code {{
    font-family: 'Cascadia Code', 'Consolas', 'Courier New', monospace;
    font-size: {cfg.code_font_size_pt}pt;
}}
p code, li code, td code, h1 code, h2 code, h3 code, h4 code, a code {{
    background: #eff1f3;
    padding: 1px 5px;
    border-radius: 3px;
    font-size: {cfg.code_font_size_pt + 1}pt;
}}

/* ── Tables — KEEP TOGETHER ────────────────────────────────── */
table {{
    border-collapse: collapse;
    width: 100%;
    margin: 14px 0;
    font-size: 10pt;
    break-inside: avoid;
    page-break-inside: avoid;
}}
thead {{ background: #1a56db; color: white; }}
th {{ padding: 8px 10px; text-align: left; font-weight: 600; }}
td {{ padding: 6px 10px; border: 1px solid #d0d7de; }}
tbody tr:nth-child(even) {{ background: #f6f8fa; }}
tbody tr:nth-child(odd) {{ background: #fff; }}

/* ── Blockquotes — KEEP TOGETHER ───────────────────────────── */
blockquote {{
    border-left: 4px solid #3b82f6;
    background: #eff6ff;
    margin: 14px 0;
    padding: 10px 16px;
    color: #1e3a5f;
    font-style: italic;
    break-inside: avoid;
    page-break-inside: avoid;
}}
blockquote p {{ margin: 4px 0; }}

/* ── Lists ─────────────────────────────────────────────────── */
ul, ol {{ margin: 8px 0; padding-left: 24px; }}
li {{ margin-bottom: 4px; }}
li > ul, li > ol {{ margin-top: 4px; margin-bottom: 4px; }}

/* ── Paragraphs & misc ─────────────────────────────────────── */
p {{ margin: 10px 0; orphans: 3; widows: 3; }}
hr {{ border: none; border-top: 2px solid #e5e7eb; margin: 22px 0; }}
a {{ color: #1a56db; text-decoration: none; }}
strong {{ color: #111; }}
h1 + *, h2 + *, h3 + *, h4 + * {{ break-before: avoid; }}

/* ── How to Use ────────────────────────────────────────────── */
.howto h1 {{ text-align: center; border-bottom: none; margin-bottom: 24px; }}
.howto h3 {{ color: #1a56db; margin-top: 22px; font-size: 13pt; }}
.howto .depth-table td {{
    padding: 4px 10px; border: none; vertical-align: top; font-size: 11pt;
}}
.howto .depth-table td:first-child {{
    font-weight: 700; color: #1a56db; white-space: nowrap;
}}

/* ── Table of Contents ─────────────────────────────────────── */
.toc h1 {{ text-align: center; border-bottom: none; margin-bottom: 28px; }}
.toc-part-title {{
    font-size: 13pt; font-weight: 700; color: #1a56db;
    margin-top: 20px; margin-bottom: 8px; padding-bottom: 4px;
    border-bottom: 1px solid #e0e0e0;
}}
.toc-entry {{
    display: flex; align-items: baseline; text-decoration: none;
    color: #1a1a1a; padding: 4px 0; font-size: 11.5pt; line-height: 1.6;
}}
.toc-entry:hover {{ color: #1a56db; }}
.toc-num {{
    display: inline-block; width: 32px; color: #1a56db;
    font-weight: 600; flex-shrink: 0;
}}
.toc-title {{ flex-shrink: 0; }}
.toc-dots {{
    flex: 1; border-bottom: 1px dotted #bbb;
    margin: 0 8px 4px 8px; min-width: 30px;
}}
.toc-page {{
    flex-shrink: 0; min-width: 26px; text-align: right;
    color: #666; font-variant-numeric: tabular-nums; font-weight: 500;
}}
"""


# ═══════════════════════════════════════════════════════════════════════
# Markdown → HTML
# ═══════════════════════════════════════════════════════════════════════

def md_to_html(text: str) -> str:
    """Convert markdown text to HTML."""
    md = markdown.Markdown(
        extensions=[TableExtension(), FencedCodeExtension(), "nl2br"],
        output_format="html5",
    )
    return md.convert(text)


def mark_long_pre_blocks(html: str) -> str:
    """Allow page breaks in code blocks longer than ~55 lines."""
    def replacer(m):
        content = m.group(1)
        if content.count("\n") > 55:
            return f'<pre class="allow-break">{content}</pre>'
        return m.group(0)
    return re.sub(r"<pre>(.*?)</pre>", replacer, html, flags=re.DOTALL)


def clean_text_for_search(title: str) -> str:
    """Strip markdown formatting for PDF text search."""
    t = re.sub(r"[`*_~#]", "", title)
    t = re.sub(r"\s+", " ", t)
    return t.strip()


def extract_h1(md_text: str) -> str:
    """Extract the first H1 heading from markdown text."""
    match = re.search(r"^#\s+(.+)$", md_text, re.MULTILINE)
    return match.group(1) if match else ""


# ═══════════════════════════════════════════════════════════════════════
# HTML Building
# ═══════════════════════════════════════════════════════════════════════

def build_cover(cfg: BookConfig) -> str:
    subtitle = f'<div class="subtitle">{cfg.subtitle}</div>' if cfg.subtitle else ""
    desc = f'<div class="tagline">{cfg.description}</div>' if cfg.description else ""
    return f"""
    <div class="cover">
        <h1>{cfg.title}</h1>
        {subtitle}
        {desc}
    </div>"""


def build_howto(cfg: BookConfig) -> str:
    """Build 'How to Use This Book' page. Uses custom text or auto-generates."""
    if cfg.howto_text:
        content = md_to_html(cfg.howto_text)
        return f'<div class="howto chapter-break"><h1>How to Use This Book</h1>{content}</div>'

    # Count total chapters
    total_chapters = sum(len(p["chapters"]) for p in cfg.parts)
    total_parts = len(cfg.parts)

    # Build parts summary
    parts_summary = ""
    for part in cfg.parts:
        ch_count = len(part["chapters"])
        parts_summary += f"<li><strong>{part['name']}</strong> ({ch_count} chapters)</li>\n"

    return f"""
    <div class="howto chapter-break">
        <h1>How to Use This Book</h1>

        <h3>Structure</h3>
        <p>This book contains <strong>{total_chapters} chapters</strong> organized
        in <strong>{total_parts} parts</strong>:</p>
        <ul>{parts_summary}</ul>

        <h3>Navigation</h3>
        <ul>
            <li><strong>Table of Contents</strong> — Click any chapter title to jump directly to it</li>
            <li><strong>Bookmark Panel</strong> — Open your PDF reader's sidebar for the full outline</li>
            <li><strong>Running Headers</strong> — The current chapter name appears at the top of every page</li>
            <li><strong>Page Numbers</strong> — Centered at the bottom, matching your reader's page indicator</li>
        </ul>

        <h3>Code &amp; Diagrams</h3>
        <ul>
            <li>Code blocks are kept on a single page whenever possible</li>
            <li>Tables and diagrams will not split across page boundaries</li>
            <li>All links in the Table of Contents are clickable</li>
        </ul>
    </div>"""


def build_toc(cfg: BookConfig, page_map: dict = None) -> str:
    """Build clickable Table of Contents. page_map: {chapter_id: page_number}."""
    entries = []
    for part in cfg.parts:
        entries.append(f'<div class="toc-part-title">{part["name"]}</div>')
        for ch in part["chapters"]:
            page_num = page_map.get(ch["id"], "...") if page_map else "..."
            entries.append(f"""
                <a href="#{ch['id']}" class="toc-entry">
                    <span class="toc-num">{ch['num']}</span>
                    <span class="toc-title">{ch['short']}</span>
                    <span class="toc-dots"></span>
                    <span class="toc-page">{page_num}</span>
                </a>""")

    return f"""
    <div class="toc chapter-break">
        <h1>Table of Contents</h1>
        {''.join(entries)}
    </div>"""


def build_full_html(cfg: BookConfig, chapters_content: list, page_map: dict = None) -> str:
    """Assemble the complete HTML document."""
    css = build_css(cfg)
    parts = [build_cover(cfg), build_howto(cfg), build_toc(cfg, page_map)]

    for ch_id, html_content in chapters_content:
        parts.append(f'<div id="{ch_id}" class="chapter-break">\n{html_content}\n</div>')

    body = "\n".join(parts)
    # Empty <title> prevents Chrome from using it as a PDF header
    return f"""<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<title> </title>
<style>{css}</style>
</head>
<body>
{body}
</body>
</html>"""


# ═══════════════════════════════════════════════════════════════════════
# File I/O & Chrome
# ═══════════════════════════════════════════════════════════════════════

def find_chrome(cfg: BookConfig) -> str:
    """Find Chrome/Edge executable."""
    if cfg.chrome_path and os.path.exists(cfg.chrome_path):
        return cfg.chrome_path
    for path in CHROME_PATHS:
        if os.path.exists(path):
            return path
    print("ERROR: Chrome or Edge not found. Specify with --chrome or in config.")
    sys.exit(1)


def read_chapters(cfg: BookConfig) -> list:
    """Read all chapter markdown files."""
    all_chapters = []
    for part in cfg.parts:
        for ch in part["chapters"]:
            path = os.path.join(cfg.source_dir, ch["file"])
            if not os.path.exists(path):
                print(f"  WARNING: File not found: {path}")
                continue

            with open(path, "r", encoding="utf-8") as f:
                md_text = f.read()

            h1 = extract_h1(md_text) or ch["short"]
            html = md_to_html(md_text)
            html = mark_long_pre_blocks(html)

            all_chapters.append({
                "id": ch["id"],
                "h1": h1,
                "h1_clean": clean_text_for_search(h1),
                "short": ch["short"],
                "num": ch["num"],
                "html": html,
            })
    return all_chapters


def chrome_to_pdf(chrome_path: str, html_path: str, pdf_path: str) -> bool:
    """Use Chrome/Edge headless to convert HTML to PDF."""
    abs_html = os.path.abspath(html_path).replace("\\", "/")
    abs_pdf = os.path.abspath(pdf_path)

    cmd = [
        chrome_path,
        "--headless=new",
        "--disable-gpu",
        "--no-sandbox",
        "--run-all-compositor-stages-before-draw",
        f"--print-to-pdf={abs_pdf}",
        "--no-pdf-header-footer",
        "--print-to-pdf-no-header",
        f"file:///{abs_html}",
    ]
    result = subprocess.run(cmd, capture_output=True, text=True, timeout=300)
    if result.returncode != 0:
        print(f"  Chrome error: {result.stderr[:500]}")
    return os.path.exists(abs_pdf)


# ═══════════════════════════════════════════════════════════════════════
# PyMuPDF Post-Processing
# ═══════════════════════════════════════════════════════════════════════

def find_chapter_pages(pdf_path: str, chapters: list) -> dict:
    """Search PDF for each chapter's H1 title to find its start page."""
    doc = fitz.open(pdf_path)

    # Build normalized text cache (collapse newlines for matching)
    page_texts = []
    for pno in range(len(doc)):
        raw = doc[pno].get_text()
        normalized = raw.replace("\n", " ").replace("  ", " ")
        page_texts.append(normalized)

    page_map = {}
    for ch in chapters:
        searches = [
            (ch["h1_clean"][:40], 4),                          # Full H1 prefix
            (re.sub(r"[^a-zA-Z0-9 ]", "", ch["h1"])[:35], 4), # ASCII-only H1
            (ch["short"], 5),                                   # Short title (skip TOC)
        ]

        found = False
        for search_text, min_page_idx in searches:
            search_text = search_text.strip()
            if not search_text or len(search_text) < 5:
                continue
            for pno in range(min(min_page_idx, len(doc)), len(doc)):
                if search_text in page_texts[pno]:
                    page_map[ch["id"]] = pno + 1  # 1-indexed
                    found = True
                    break
            if found:
                break

        if not found:
            print(f"  WARNING: Could not find chapter {ch['id']} ('{ch['short']}')")

    doc.close()
    return page_map


def build_page_to_chapter(page_map: dict, total_pages: int, chapters: list) -> dict:
    """Map every page to its chapter for running headers."""
    sorted_chs = sorted(
        [(ch["id"], ch["short"], ch["num"], page_map.get(ch["id"], 9999))
         for ch in chapters],
        key=lambda x: x[3],
    )

    mapping = {}
    for i, (ch_id, short, num, start) in enumerate(sorted_chs):
        end = sorted_chs[i + 1][3] - 1 if i + 1 < len(sorted_chs) else total_pages
        for p in range(start, end + 1):
            mapping[p] = (num, short)

    return mapping


def post_process(input_pdf: str, output_pdf: str, cfg: BookConfig,
                 chapters: list, page_map: dict):
    """Add bookmarks, running headers, page numbers, and metadata."""
    doc = fitz.open(input_pdf)
    total = len(doc)
    print(f"  Processing {total} pages...")

    first_chapter_page = min(page_map.values()) if page_map else 4
    pg_ch_map = build_page_to_chapter(page_map, total, chapters)

    # Placement constants
    header_font = "helv"
    header_y = cfg.header_text_y_mm * MM
    line_y = cfg.header_line_y_mm * MM
    footer_y = cfg.footer_text_y_mm * MM
    margin_x = cfg.margin_lr_mm * MM
    header_color = (0.35, 0.35, 0.35)
    line_color = (0.78, 0.78, 0.78)
    footer_color = (0.4, 0.4, 0.4)

    for pno in range(total):
        page = doc[pno]
        page_num = pno + 1

        if page_num == 1:
            continue  # No header/footer on cover

        # ── Page number (bottom center) ──
        num_text = str(page_num)
        tw = fitz.get_text_length(num_text, fontname=header_font,
                                  fontsize=cfg.footer_font_size_pt)
        page.insert_text(
            fitz.Point((PAGE_W - tw) / 2, footer_y),
            num_text, fontname=header_font,
            fontsize=cfg.footer_font_size_pt, color=footer_color,
        )

        # ── Running header (content pages, not chapter starts) ──
        if page_num >= first_chapter_page and page_num in pg_ch_map:
            is_chapter_start = page_num in page_map.values()
            if not is_chapter_start:
                ch_num, ch_short = pg_ch_map[page_num]
                header_text = f"{ch_num} - {ch_short}"

                page.insert_text(
                    fitz.Point(margin_x, header_y),
                    header_text, fontname=header_font,
                    fontsize=cfg.header_font_size_pt, color=header_color,
                )
                page.draw_line(
                    fitz.Point(margin_x, line_y),
                    fitz.Point(PAGE_W - margin_x, line_y),
                    color=line_color, width=0.5,
                )

    # ── Bookmarks ──
    howto_page = 2
    toc_page = 3
    for pno in range(min(6, total)):
        text = doc[pno].get_text()
        if "How to Use This Book" in text and pno + 1 != 1:
            howto_page = pno + 1
        if "Table of Contents" in text and pno + 1 > 2:
            toc_page = pno + 1

    toc_entries = [
        [1, "Cover", 1],
        [1, "How to Use This Book", howto_page],
        [1, "Table of Contents", toc_page],
    ]
    for part in cfg.parts:
        first_ch_id = part["chapters"][0]["id"]
        part_page = page_map.get(first_ch_id, 1)
        toc_entries.append([1, part["name"], part_page])
        for ch in part["chapters"]:
            ch_page = page_map.get(ch["id"], 1)
            toc_entries.append([2, f'{ch["num"]} - {ch["short"]}', ch_page])

    doc.set_toc(toc_entries)
    print(f"  Added {len(toc_entries)} bookmarks")

    # ── Metadata ──
    doc.set_metadata({
        "title": f"{cfg.title}" + (f" - {cfg.subtitle}" if cfg.subtitle else ""),
        "author": cfg.author,
        "subject": cfg.subject or cfg.description,
        "keywords": cfg.keywords,
        "creator": "md2pdf (Chrome Headless + PyMuPDF)",
        "producer": "md2pdf",
    })

    # ── Page labels ──
    doc.set_page_labels([{
        "startpage": 0, "prefix": "", "style": "D", "firstpagenum": 1
    }])

    # ── Save ──
    doc.save(output_pdf, garbage=4, deflate=True, clean=True)
    doc.close()

    size_mb = os.path.getsize(output_pdf) / 1024 / 1024
    print(f"  Saved: {output_pdf} ({size_mb:.1f} MB, {total} pages)")


# ═══════════════════════════════════════════════════════════════════════
# Main Pipeline
# ═══════════════════════════════════════════════════════════════════════

def generate(cfg: BookConfig):
    """Full PDF generation pipeline."""
    print("=" * 60)
    print(f"  md2pdf: {cfg.title}")
    print("=" * 60)

    chrome = find_chrome(cfg)
    cfg.auto_detect_chapters()
    cfg.resolve_chapters()

    # ── Read chapters ──
    print("\n[1/6] Reading chapters...")
    chapters = read_chapters(cfg)
    print(f"  Found {len(chapters)} chapters")
    if not chapters:
        print("ERROR: No chapters found.")
        return 1

    chapters_content = [(ch["id"], ch["html"]) for ch in chapters]

    temp_html = os.path.join(os.path.dirname(cfg.output), "_md2pdf_temp.html")
    temp_pdf = os.path.join(os.path.dirname(cfg.output), "_md2pdf_temp.pdf")
    chrome_pdf = os.path.join(os.path.dirname(cfg.output), "_md2pdf_chrome.pdf")

    # ── Pass 1: temp PDF for page discovery ──
    print("\n[2/6] Pass 1 — temp PDF for page discovery...")
    html1 = build_full_html(cfg, chapters_content, page_map=None)
    with open(temp_html, "w", encoding="utf-8") as f:
        f.write(html1)
    if not chrome_to_pdf(chrome, temp_html, temp_pdf):
        print("  ERROR: Chrome failed on pass 1")
        return 1

    # ── Extract page numbers ──
    print("\n[3/6] Extracting chapter page numbers...")
    page_map = find_chapter_pages(temp_pdf, chapters)
    found = len(page_map)
    print(f"  Found {found}/{len(chapters)} chapters")
    for ch in chapters:
        pg = page_map.get(ch["id"], "?")
        print(f"    {ch['num']} - {ch['short']}: page {pg}")

    # ── Pass 2: final PDF with real page numbers ──
    print("\n[4/6] Pass 2 — PDF with TOC page numbers...")
    html2 = build_full_html(cfg, chapters_content, page_map=page_map)
    with open(temp_html, "w", encoding="utf-8") as f:
        f.write(html2)
    if not chrome_to_pdf(chrome, temp_html, chrome_pdf):
        print("  ERROR: Chrome failed on pass 2")
        return 1

    # Verify stability
    page_map2 = find_chapter_pages(chrome_pdf, chapters)
    shifts = sum(1 for cid in page_map if page_map.get(cid) != page_map2.get(cid))
    if shifts:
        print(f"  {shifts} chapters shifted — pass 3...")
        html3 = build_full_html(cfg, chapters_content, page_map=page_map2)
        with open(temp_html, "w", encoding="utf-8") as f:
            f.write(html3)
        chrome_to_pdf(chrome, temp_html, chrome_pdf)
        page_map = page_map2
    else:
        print("  Page numbers stable")

    # ── Post-process ──
    print("\n[5/6] Post-processing...")
    post_process(chrome_pdf, cfg.output, cfg, chapters, page_map)

    # ── Cleanup ──
    print("\n[6/6] Cleaning up...")
    for path in [temp_html, temp_pdf, chrome_pdf]:
        if os.path.exists(path):
            os.remove(path)

    final_size = os.path.getsize(cfg.output) / 1024 / 1024
    doc = fitz.open(cfg.output)
    final_pages = len(doc)
    doc.close()

    print(f"\n{'=' * 60}")
    print(f"  Done: {cfg.output}")
    print(f"  {final_pages} pages, {final_size:.1f} MB")
    print(f"  Clickable TOC, bookmarks, running headers, page numbers")
    print(f"{'=' * 60}")
    return 0


# ═══════════════════════════════════════════════════════════════════════
# CLI
# ═══════════════════════════════════════════════════════════════════════

def main():
    parser = argparse.ArgumentParser(
        prog="md2pdf",
        description="Convert a collection of Markdown files into a professional PDF book.",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  # From config file (recommended):
  python tools/md2pdf.py --config tools/book_config.json

  # Quick auto-detect (numbered .md files):
  python tools/md2pdf.py --source-dir docs/ --title "My Notes" --output notes.pdf

  # With custom Chrome path:
  python tools/md2pdf.py --config book.json --chrome /usr/bin/chromium
        """,
    )
    parser.add_argument("--config", "-c", help="Path to JSON config file")
    parser.add_argument("--source-dir", "-s", default=".", help="Directory containing .md files")
    parser.add_argument("--output", "-o", default="book.pdf", help="Output PDF path")
    parser.add_argument("--title", "-t", help="Book title")
    parser.add_argument("--subtitle", default="", help="Book subtitle")
    parser.add_argument("--author", "-a", default="", help="Author name")
    parser.add_argument("--chrome", help="Path to Chrome/Edge executable")

    args = parser.parse_args()

    if args.config:
        cfg = BookConfig.from_json(args.config)
        # CLI overrides
        if args.chrome:
            cfg.chrome_path = args.chrome
        if args.title:
            cfg.title = args.title
        if args.output != "book.pdf":
            cfg.output = os.path.abspath(args.output)
    else:
        cfg = BookConfig.from_args(args)

    return generate(cfg)


if __name__ == "__main__":
    sys.exit(main())
