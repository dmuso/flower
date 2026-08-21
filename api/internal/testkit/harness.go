package testkit

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"flower/api/internal/app"
	clk "flower/api/internal/platform/clock"
	"flower/api/internal/platform/config"
	"flower/api/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Harness struct {
	T      *testing.T
	App    *app.App
	DB     *sql.DB
	Now    time.Time
	Origin string
}

var dbMu sync.Mutex

func lockTestDB(t *testing.T) {
	t.Helper()
	// Packages run as separate binaries; a process-local mutex is not enough.
	f, err := os.OpenFile("/tmp/flower-api-test.lock", os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("test lock: %v", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatalf("test flock: %v", err)
	}
	dbMu.Lock()
	t.Cleanup(func() {
		dbMu.Unlock()
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	})
}

func New(t *testing.T) *Harness {
	t.Helper()
	lockTestDB(t)
	gin.SetMode(gin.TestMode)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("testkit: cannot locate caller")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", "..", ".."))
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir repo root: %v", err)
	}

	if os.Getenv("ENVIRONMENT") == "" {
		t.Setenv("ENVIRONMENT", "test")
	}
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config: %v", err)
	}
	now := time.Date(2026, 8, 21, 0, 10, 0, 0, time.UTC)
	application, err := app.New(cfg, app.WithClock(clk.Fixed{T: now}))
	if err != nil {
		t.Fatalf("app.New: %v (is Postgres on %s:%s ?)", err, cfg.Database.Host, cfg.Database.Port)
	}
	t.Cleanup(func() {
		if err := application.Close(); err != nil && !strings.Contains(err.Error(), "invalid argument") {
			t.Logf("close app: %v", err)
		}
	})
	h := &Harness{T: t, App: application, DB: application.DB, Now: now, Origin: cfg.FrontendOrigin}
	h.Reset()
	return h
}

func (h *Harness) Reset() {
	h.T.Helper()
	_, err := h.DB.Exec(`
		TRUNCATE
			email_outbox, auth_tokens, sessions,
			story_labels, labels, stories, activities, iterations,
			project_memberships, organisation_memberships,
			projects, organisations, users
		RESTART IDENTITY CASCADE
	`)
	if err != nil {
		h.T.Fatalf("truncate: %v", err)
	}
}

type Response struct {
	Code    int
	Body    []byte
	Header  http.Header
	Cookies []*http.Cookie
}

func (r *Response) JSON() map[string]any {
	var m map[string]any
	if err := json.Unmarshal(r.Body, &m); err != nil {
		return map[string]any{"_raw": string(r.Body)}
	}
	return m
}

func (h *Harness) Do(method, path string, body any, cookies ...*http.Cookie) *Response {
	h.T.Helper()
	var reader io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			h.T.Fatalf("marshal: %v", err)
		}
		reader = bytes.NewReader(raw)
	}
	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", h.Origin)
	for _, c := range cookies {
		req.AddCookie(c)
	}
	rec := httptest.NewRecorder()
	h.App.Router.ServeHTTP(rec, req)
	return &Response{
		Code:    rec.Code,
		Body:    rec.Body.Bytes(),
		Header:  rec.Header(),
		Cookies: rec.Result().Cookies(),
	}
}

func (h *Harness) SessionCookie(res *Response) *http.Cookie {
	h.T.Helper()
	for _, c := range res.Cookies {
		if c.Name == middleware.CookieName {
			return c
		}
	}
	h.T.Fatalf("missing %s cookie; status=%d body=%s", middleware.CookieName, res.Code, res.Body)
	return nil
}

func (h *Harness) OutboxToken(email, kind string) string {
	h.T.Helper()
	var body string
	err := h.DB.QueryRow(`
		SELECT body FROM email_outbox WHERE to_email = $1 AND kind = $2 ORDER BY created_at DESC LIMIT 1
	`, email, kind).Scan(&body)
	if err != nil {
		h.T.Fatalf("outbox %s/%s: %v", email, kind, err)
	}
	re := regexp.MustCompile(`token=([0-9a-fA-F]+)`)
	m := re.FindStringSubmatch(body)
	if len(m) != 2 {
		h.T.Fatalf("no token in outbox body: %s", body)
	}
	return m[1]
}

func (h *Harness) Count(table string) int {
	h.T.Helper()
	var n int
	if err := h.DB.QueryRow("SELECT COUNT(1) FROM " + table).Scan(&n); err != nil {
		h.T.Fatalf("count %s: %v", table, err)
	}
	return n
}

func (h *Harness) SignUp(email, password string) *Response {
	return h.Do(http.MethodPost, "/api/v1/auth/signup", map[string]string{"email": email, "password": password})
}

func (h *Harness) Verify(email string) *http.Cookie {
	h.T.Helper()
	token := h.OutboxToken(email, "verify_email")
	res := h.Do(http.MethodPost, "/api/v1/auth/verify-email/consume", map[string]string{"token": token})
	if res.Code != http.StatusOK {
		h.T.Fatalf("verify: status %d body %s", res.Code, res.Body)
	}
	return h.SessionCookie(res)
}

func ErrorCode(res *Response) string {
	m := res.JSON()
	errObj, _ := m["error"].(map[string]any)
	if errObj == nil {
		return ""
	}
	code, _ := errObj["code"].(string)
	return code
}

func StringField(m map[string]any, key string) string {
	v, _ := m[key].(string)
	return v
}

func MustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q in %q", needle, haystack)
	}
}
