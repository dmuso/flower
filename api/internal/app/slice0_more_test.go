package app_test

import (
	"net/http"
	"strings"
	"testing"

	"flower/api/internal/testkit"
)

func TestMeRequiresSession(t *testing.T) {
	h := testkit.New(t)
	res := h.Do(http.MethodGet, "/api/v1/me", nil)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("status %d", res.Code)
	}
}

func TestConsumedVerifyLinkFails(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("reuse@example.com", "secret12")
	token := h.OutboxToken("reuse@example.com", "verify_email")
	first := h.Do(http.MethodPost, "/api/v1/auth/verify-email/consume", map[string]string{"token": token})
	if first.Code != http.StatusOK {
		t.Fatalf("first %d %s", first.Code, first.Body)
	}
	second := h.Do(http.MethodPost, "/api/v1/auth/verify-email/consume", map[string]string{"token": token})
	if second.Code != http.StatusUnauthorized {
		t.Fatalf("second %d %s", second.Code, second.Body)
	}
}

func TestLogoutClearsSession(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	cookie := h.Verify("maya@example.com")
	out := h.Do(http.MethodPost, "/api/v1/auth/logout", nil, cookie)
	if out.Code != http.StatusNoContent {
		t.Fatalf("logout %d", out.Code)
	}
	me := h.Do(http.MethodGet, "/api/v1/me", nil, cookie)
	if me.Code != http.StatusUnauthorized {
		t.Fatalf("me after logout %d", me.Code)
	}
}

func TestSessionCookieFlags(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	token := h.OutboxToken("maya@example.com", "verify_email")
	res := h.Do(http.MethodPost, "/api/v1/auth/verify-email/consume", map[string]string{"token": token})
	if res.Code != http.StatusOK {
		t.Fatalf("verify %d %s", res.Code, res.Body)
	}
	cookie := h.SessionCookie(res)
	if !cookie.HttpOnly {
		t.Fatal("cookie must be HttpOnly")
	}
	if cookie.Name != "flower_session" {
		t.Fatalf("name %s", cookie.Name)
	}
	raw := strings.Join(res.Header.Values("Set-Cookie"), ";")
	if !strings.Contains(strings.ToLower(raw), "samesite=lax") {
		t.Fatalf("SameSite=Lax missing from %s", raw)
	}
}

func TestListProjectsForOrganisation(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	cookie := h.Verify("maya@example.com")
	orgID := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie).JSON(), "id")
	h.Do(http.MethodPost, "/api/v1/organisations/"+orgID+"/projects", map[string]string{"name": "Trail"}, cookie)
	list := h.Do(http.MethodGet, "/api/v1/organisations/"+orgID+"/projects", nil, cookie)
	if list.Code != http.StatusOK {
		t.Fatalf("list %d %s", list.Code, list.Body)
	}
	projects, _ := list.JSON()["projects"].([]any)
	if len(projects) != 1 {
		t.Fatalf("projects %v", list.JSON())
	}
}
