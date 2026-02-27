@echo off
cd /d C:\Users\samte\GolandProjects\go-interview-prep
echo === Git Status ===
git status
echo === Adding all changes ===
git add -A
echo === Committing ===
git commit -m "test: rewrite all exercise test files with consistent pass/fail output style"
echo === Pushing to main ===
git push origin main
echo === Done ===
pause

