package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: TimingMiddleware
// ────────────────────────────────────────────────────────────

func TestTimingMiddleware(t *testing.T) {
	var recorded time.Duration
	mw := TimingMiddleware(func(d time.Duration) { recorded = d })

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.Write([]byte("ok"))
	})

	handler := mw(inner)
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if recorded < 10*time.Millisecond {
		t.Errorf("❌ TimingMiddleware recorded %v, want >= 10ms\n\t\t"+
			"Hint: start := time.Now(); next.ServeHTTP(w, r); onDone(time.Since(start))",
			recorded)
	} else {
		t.Logf("✅ TimingMiddleware recorded %v", recorded)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: CORSMiddleware
// ────────────────────────────────────────────────────────────

func TestCORSMiddleware(t *testing.T) {
	mw := CORSMiddleware("https://example.com")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// Normal request
	req := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Error("❌ CORS origin not set\n\t\t" +
			"Hint: w.Header().Set(\"Access-Control-Allow-Origin\", allowedOrigin)")
	}

	// OPTIONS preflight
	req = httptest.NewRequest("OPTIONS", "/api", nil)
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("❌ OPTIONS status = %d, want 204", w.Code)
	}
	if w.Body.String() != "" {
		t.Error("❌ OPTIONS should not call next handler")
	} else {
		t.Log("✅ CORSMiddleware handles preflight + regular requests")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: RequestIDMiddleware
// ────────────────────────────────────────────────────────────

func TestRequestIDMiddleware(t *testing.T) {
	mw := RequestIDMiddleware(func() string { return "generated-id" })
	var capturedID string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID, _ = r.Context().Value(RequestIDKey).(string)
		w.Write([]byte("ok"))
	})

	// Without existing header (should generate)
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "generated-id" {
		t.Errorf("❌ response X-Request-ID = %q, want \"generated-id\"\n\t\t"+
			"Hint: id := r.Header.Get(\"X-Request-ID\"); if id == \"\" { id = generateID() }; "+
			"w.Header().Set(\"X-Request-ID\", id); ctx := context.WithValue(r.Context(), RequestIDKey, id)",
			w.Header().Get("X-Request-ID"))
	}
	if capturedID != "generated-id" {
		t.Errorf("❌ context ID = %q, want \"generated-id\"", capturedID)
	} else {
		t.Log("✅ RequestIDMiddleware generates and propagates ID")
	}

	// With existing header (should reuse)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") != "existing-id" {
		t.Error("❌ should reuse existing X-Request-ID")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: ContentTypeMiddleware
// ────────────────────────────────────────────────────────────

func TestContentTypeMiddleware(t *testing.T) {
	mw := ContentTypeMiddleware("application/json")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// GET should pass through
	req := httptest.NewRequest("GET", "/api", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Body.String() != "ok" {
		t.Error("❌ GET should pass through without checking Content-Type")
	}

	// POST with wrong Content-Type
	req = httptest.NewRequest("POST", "/api", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "text/plain")
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Code != 415 {
		t.Errorf("❌ POST with wrong Content-Type: status = %d, want 415\n\t\t"+
			"Hint: for POST/PUT/PATCH, check r.Header.Get(\"Content-Type\")",
			w.Code)
	} else {
		t.Log("✅ ContentTypeMiddleware rejects wrong Content-Type")
	}

	// POST with correct Content-Type
	req = httptest.NewRequest("POST", "/api", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Body.String() != "ok" {
		t.Error("❌ POST with correct Content-Type should pass through")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: MaxBytesMiddleware
// ────────────────────────────────────────────────────────────

func TestMaxBytesMiddleware(t *testing.T) {
	mw := MaxBytesMiddleware(10)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "too large", http.StatusRequestEntityTooLarge)
			return
		}
		w.Write([]byte("ok"))
	})

	// Small body
	req := httptest.NewRequest("POST", "/", strings.NewReader("small"))
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Body.String() != "ok" {
		t.Error("❌ small body should pass through\n\t\t" +
			"Hint: r.Body = http.MaxBytesReader(w, r.Body, maxBytes)")
	}

	// Large body
	req = httptest.NewRequest("POST", "/", strings.NewReader("this is way too large"))
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("❌ large body status = %d, want 413", w.Code)
	} else {
		t.Log("✅ MaxBytesMiddleware limits body size")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: BasicAuthMiddleware
// ────────────────────────────────────────────────────────────

func TestBasicAuthMiddleware(t *testing.T) {
	mw := BasicAuthMiddleware("admin", "secret")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("protected"))
	})

	// With correct credentials
	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("admin", "secret")
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Body.String() != "protected" {
		t.Error("❌ correct credentials should pass through\n\t\t" +
			"Hint: user, pass, ok := r.BasicAuth(); if !ok || user != username || pass != password { " +
			"w.Header().Set(\"WWW-Authenticate\", `Basic realm=\"restricted\"`); w.WriteHeader(401) }")
	} else {
		t.Log("✅ BasicAuth with correct credentials passes")
	}

	// Without credentials
	req = httptest.NewRequest("GET", "/", nil)
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Code != 401 {
		t.Errorf("❌ no credentials: status = %d, want 401", w.Code)
	} else {
		t.Log("✅ BasicAuth rejects missing credentials")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: CacheControlMiddleware
// ────────────────────────────────────────────────────────────

func TestCacheControlMiddleware(t *testing.T) {
	mw := CacheControlMiddleware("max-age=3600")
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("cached"))
	})

	req := httptest.NewRequest("GET", "/static", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Header().Get("Cache-Control") != "max-age=3600" {
		t.Errorf("❌ Cache-Control = %q, want \"max-age=3600\"", w.Header().Get("Cache-Control"))
	} else {
		t.Log("✅ CacheControlMiddleware sets header")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: SecurityHeadersMiddleware
// ────────────────────────────────────────────────────────────

func TestSecurityHeadersMiddleware(t *testing.T) {
	mw := SecurityHeadersMiddleware()
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secure"))
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	expected := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"X-XSS-Protection":          "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
	}

	for header, want := range expected {
		got := w.Header().Get(header)
		if got != want {
			t.Errorf("❌ %s = %q, want %q", header, got, want)
		}
	}
	t.Log("✅ SecurityHeadersMiddleware sets all headers")
}

// ────────────────────────────────────────────────────────────
// Exercise 9: TimeoutMiddleware
// ────────────────────────────────────────────────────────────

func TestTimeoutMiddleware(t *testing.T) {
	mw := TimeoutMiddleware(50 * time.Millisecond)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			return
		case <-time.After(200 * time.Millisecond):
			w.Write([]byte("too slow"))
		}
	})

	req := httptest.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if w.Body.String() == "too slow" {
		t.Error("❌ handler should have been cancelled by timeout\n\t\t" +
			"Hint: use http.TimeoutHandler(next, timeout, \"timeout\") or " +
			"ctx, cancel := context.WithTimeout(r.Context(), timeout)")
	} else {
		t.Log("✅ TimeoutMiddleware cancels slow handler")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: MethodOverrideMiddleware
// ────────────────────────────────────────────────────────────

func TestMethodOverrideMiddleware(t *testing.T) {
	mw := MethodOverrideMiddleware()
	var capturedMethod string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedMethod = r.Method
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest("POST", "/resource", nil)
	req.Header.Set("X-HTTP-Method-Override", "PUT")
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)

	if capturedMethod != "PUT" {
		t.Errorf("❌ method = %q, want \"PUT\"\n\t\t"+
			"Hint: if override := r.Header.Get(\"X-HTTP-Method-Override\"); override != \"\" { r.Method = override }",
			capturedMethod)
	} else {
		t.Log("✅ MethodOverrideMiddleware changes POST → PUT")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ConditionalMiddleware
// ────────────────────────────────────────────────────────────

func TestConditionalMiddleware(t *testing.T) {
	addHeader := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Conditional", "applied")
			next.ServeHTTP(w, r)
		})
	}
	onlyGET := func(r *http.Request) bool { return r.Method == "GET" }
	mw := ConditionalMiddleware(onlyGET, addHeader)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// GET — middleware should apply
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Header().Get("X-Conditional") != "applied" {
		t.Error("❌ GET should have conditional middleware applied\n\t\t" +
			"Hint: if condition(r) { mw(next).ServeHTTP(w, r) } else { next.ServeHTTP(w, r) }")
	}

	// POST — middleware should not apply
	req = httptest.NewRequest("POST", "/", nil)
	w = httptest.NewRecorder()
	mw(inner).ServeHTTP(w, req)
	if w.Header().Get("X-Conditional") == "applied" {
		t.Error("❌ POST should NOT have conditional middleware applied")
	} else {
		t.Log("✅ ConditionalMiddleware applies only when condition is true")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: Chain
// ────────────────────────────────────────────────────────────

func TestChain(t *testing.T) {
	var order []string

	mw := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	handler := Chain(inner, mw("A"), mw("B"), mw("C"))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	got := strings.Join(order, ",")
	if got != "A,B,C,handler" {
		t.Errorf("❌ order = %q, want \"A,B,C,handler\"\n\t\t"+
			"Hint: iterate in reverse: for i := len(mws)-1; i >= 0; i-- { h = mws[i](h) }",
			got)
	} else {
		t.Logf("✅ Chain order: %s", got)
	}
}
