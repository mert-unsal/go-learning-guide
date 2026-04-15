package net_http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ────────────────────────────────────────────────────────────
// Exercise 1: HelloHandler
// ────────────────────────────────────────────────────────────

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	HelloHandler{}.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		t.Errorf("❌ status = %d, want 200\n\t\t"+
			"Hint: implement ServeHTTP on HelloHandler. "+
			"w.WriteHeader(200) then w.Write([]byte(\"Hello, World!\")). "+
			"Note: WriteHeader is optional for 200 (implicit)",
			resp.StatusCode)
	}
	if string(body) != "Hello, World!" {
		t.Errorf("❌ body = %q, want \"Hello, World!\"", string(body))
	} else {
		t.Logf("✅ HelloHandler responds correctly")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 2: GreetHandler
// ────────────────────────────────────────────────────────────

func TestGreetHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"with name", "/greet?name=Go", "Hello, Go!"},
		{"without name", "/greet", "Hello, stranger!"},
		{"empty name", "/greet?name=", "Hello, stranger!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			GreetHandler(w, req)

			body, _ := io.ReadAll(w.Result().Body)
			if string(body) != tt.want {
				t.Errorf("❌ GreetHandler(%s) body = %q, want %q\n\t\t"+
					"Hint: r.URL.Query().Get(\"name\") returns \"\" if missing. "+
					"Use fmt.Fprintf(w, \"Hello, %%s!\", name)",
					tt.url, string(body), tt.want)
			} else {
				t.Logf("✅ GreetHandler(%s) = %q", tt.url, tt.want)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 3: EchoHandler
// ────────────────────────────────────────────────────────────

func TestEchoHandler(t *testing.T) {
	body := "echo this back please"
	req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(body))
	w := httptest.NewRecorder()

	EchoHandler(w, req)

	resp, _ := io.ReadAll(w.Result().Body)
	if string(resp) != body {
		t.Errorf("❌ EchoHandler body = %q, want %q\n\t\t"+
			"Hint: data, _ := io.ReadAll(r.Body); defer r.Body.Close(); w.Write(data)",
			string(resp), body)
	} else {
		t.Logf("✅ EchoHandler echoed %d bytes", len(resp))
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 4: MethodRouter
// ────────────────────────────────────────────────────────────

func TestMethodRouter(t *testing.T) {
	tests := []struct {
		method     string
		wantStatus int
		wantBody   string
	}{
		{http.MethodGet, 200, "GET OK"},
		{http.MethodPost, 200, "POST OK"},
		{http.MethodDelete, 405, ""},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/route", nil)
			w := httptest.NewRecorder()
			MethodRouter(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("❌ %s status = %d, want %d\n\t\t"+
					"Hint: switch r.Method { case http.MethodGet: ... default: w.WriteHeader(405) }",
					tt.method, w.Code, tt.wantStatus)
				return
			}
			body := w.Body.String()
			if tt.wantBody != "" && body != tt.wantBody {
				t.Errorf("❌ %s body = %q, want %q", tt.method, body, tt.wantBody)
			} else {
				t.Logf("✅ MethodRouter(%s) = %d", tt.method, w.Code)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 5: JSONResponse
// ────────────────────────────────────────────────────────────

func TestJSONResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	w := httptest.NewRecorder()
	JSONResponse(w, req)

	resp := w.Result()
	ct := resp.Header.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("❌ Content-Type = %q, want \"application/json\"\n\t\t"+
			"Hint: w.Header().Set(\"Content-Type\", \"application/json\") "+
			"MUST be called BEFORE w.Write() or w.WriteHeader()",
			ct)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != `{"status":"ok"}` {
		t.Errorf("❌ body = %q, want `{\"status\":\"ok\"}`", string(body))
	} else {
		t.Logf("✅ JSONResponse: Content-Type + body correct")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 6: StatusHandler
// ────────────────────────────────────────────────────────────

func TestStatusHandler(t *testing.T) {
	tests := []struct {
		query      string
		wantStatus int
		wantBody   string
	}{
		{"?code=404", 404, "Not Found"},
		{"?code=201", 201, "Created"},
		{"?code=", 200, "OK"},
		{"", 200, "OK"},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/status"+tt.query, nil)
			w := httptest.NewRecorder()
			StatusHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("❌ status = %d, want %d\n\t\t"+
					"Hint: strconv.Atoi(r.URL.Query().Get(\"code\")); "+
					"http.StatusText(code) maps code→text",
					w.Code, tt.wantStatus)
			}
			body := w.Body.String()
			if body != tt.wantBody {
				t.Errorf("❌ body = %q, want %q", body, tt.wantBody)
			} else {
				t.Logf("✅ StatusHandler(%s) = %d %q", tt.query, w.Code, tt.wantBody)
			}
		})
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 7: SetupMux
// ────────────────────────────────────────────────────────────

func TestSetupMux(t *testing.T) {
	mux := SetupMux()

	// Test /hello
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body, _ := io.ReadAll(w.Result().Body)
	if string(body) != "Hello, World!" {
		t.Errorf("❌ GET /hello body = %q, want \"Hello, World!\"\n\t\t"+
			"Hint: mux.Handle(\"/hello\", HelloHandler{}) registers the handler. "+
			"Use mux.HandleFunc for plain functions",
			string(body))
	} else {
		t.Logf("✅ GET /hello works")
	}

	// Test /greet
	req = httptest.NewRequest(http.MethodGet, "/greet?name=Mux", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != "Hello, Mux!" {
		t.Errorf("❌ GET /greet?name=Mux = %q, want \"Hello, Mux!\"", string(body))
	} else {
		t.Logf("✅ GET /greet?name=Mux works")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 8: LoggingMiddleware
// ────────────────────────────────────────────────────────────

func TestLoggingMiddleware(t *testing.T) {
	var loggedMethod, loggedPath string
	logger := func(method, path string) {
		loggedMethod = method
		loggedPath = path
	}

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	handler := LoggingMiddleware(inner, logger)

	req := httptest.NewRequest(http.MethodPost, "/test/path", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if loggedMethod != "POST" || loggedPath != "/test/path" {
		t.Errorf("❌ logged (%q, %q), want (\"POST\", \"/test/path\")\n\t\t"+
			"Hint: return http.HandlerFunc(func(w, r) { "+
			"logger(r.Method, r.URL.Path); next.ServeHTTP(w, r) })",
			loggedMethod, loggedPath)
	} else {
		t.Logf("✅ LoggingMiddleware logged POST /test/path")
	}

	if w.Body.String() != "ok" {
		t.Errorf("❌ inner handler body = %q, want \"ok\"", w.Body.String())
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 9: RecoveryMiddleware
// ────────────────────────────────────────────────────────────

func TestRecoveryMiddleware(t *testing.T) {
	var recovered interface{}
	onPanic := func(v interface{}) { recovered = v }

	panicker := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("something went wrong")
	})

	handler := RecoveryMiddleware(panicker, onPanic)

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	// Wrap in our own recover to survive if the stub doesn't catch the panic
	didPanic := false
	func() {
		defer func() {
			if r := recover(); r != nil {
				didPanic = true
			}
		}()
		handler.ServeHTTP(w, req)
	}()

	if didPanic {
		t.Error("❌ panic propagated through RecoveryMiddleware\n\t\t" +
			"Hint: defer func() { if r := recover(); r != nil { onPanic(r); w.WriteHeader(500) } }() " +
			"BEFORE calling next.ServeHTTP")
		return
	}

	if w.Code != 500 {
		t.Errorf("❌ status = %d after panic, want 500", w.Code)
	}

	if recovered == nil || recovered != "something went wrong" {
		t.Errorf("❌ recovered = %v, want \"something went wrong\"", recovered)
	} else {
		t.Logf("✅ RecoveryMiddleware caught panic and returned 500")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 10: AuthMiddleware
// ────────────────────────────────────────────────────────────

func TestAuthMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secret data"))
	})

	handler := AuthMiddleware(inner, "Bearer my-token")

	// With correct token
	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	req.Header.Set("Authorization", "Bearer my-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != 200 || w.Body.String() != "secret data" {
		t.Errorf("❌ with valid token: status=%d body=%q\n\t\t"+
			"Hint: if r.Header.Get(\"Authorization\") == expectedToken { next.ServeHTTP(w, r) } "+
			"else { w.WriteHeader(401) }",
			w.Code, w.Body.String())
	} else {
		t.Logf("✅ AuthMiddleware allows valid token")
	}

	// Without token
	req = httptest.NewRequest(http.MethodGet, "/secret", nil)
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("❌ without token: status=%d, want 401", w.Code)
	} else {
		t.Logf("✅ AuthMiddleware rejects missing token with 401")
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 11: ChainMiddleware
// ────────────────────────────────────────────────────────────

func TestChainMiddleware(t *testing.T) {
	var order []string

	mw := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name+"-before")
				next.ServeHTTP(w, r)
				order = append(order, name+"-after")
			})
		}
	}

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
		w.Write([]byte("done"))
	})

	handler := ChainMiddleware(inner, mw("A"), mw("B"), mw("C"))

	req := httptest.NewRequest(http.MethodGet, "/chain", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Expected order: A wraps B wraps C wraps handler
	// Execution: A-before → B-before → C-before → handler → C-after → B-after → A-after
	expected := "A-before,B-before,C-before,handler,C-after,B-after,A-after"
	got := strings.Join(order, ",")
	if got != expected {
		t.Errorf("❌ execution order = %q\n\twant %q\n\t\t"+
			"Hint: iterate middlewares in REVERSE. "+
			"for i := len(middlewares)-1; i >= 0; i-- { handler = middlewares[i](handler) }",
			got, expected)
	} else {
		t.Logf("✅ ChainMiddleware onion order: %s", got)
	}
}

// ────────────────────────────────────────────────────────────
// Exercise 12: HeadersHandler
// ────────────────────────────────────────────────────────────

func TestHeadersHandler(t *testing.T) {
	// With X-Request-ID
	req := httptest.NewRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("X-Request-ID", "abc-123")
	w := httptest.NewRecorder()
	HeadersHandler(w, req)

	resp := w.Result()
	if resp.Header.Get("X-Request-ID") != "abc-123" {
		t.Errorf("❌ response X-Request-ID = %q, want \"abc-123\"\n\t\t"+
			"Hint: w.Header().Set(\"X-Request-ID\", id) BEFORE w.Write(). "+
			"Headers set after Write are ignored (superheader bug)",
			resp.Header.Get("X-Request-ID"))
	}
	if resp.Header.Get("X-Powered-By") != "go-learning-guide" {
		t.Errorf("❌ X-Powered-By = %q, want \"go-learning-guide\"",
			resp.Header.Get("X-Powered-By"))
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "abc-123" {
		t.Errorf("❌ body = %q, want \"abc-123\"", string(body))
	} else {
		t.Logf("✅ HeadersHandler echoes X-Request-ID")
	}

	// Without X-Request-ID
	req = httptest.NewRequest(http.MethodGet, "/headers", nil)
	w = httptest.NewRecorder()
	HeadersHandler(w, req)

	body, _ = io.ReadAll(w.Result().Body)
	if string(body) != "none" {
		t.Errorf("❌ body without header = %q, want \"none\"", string(body))
	}
}
