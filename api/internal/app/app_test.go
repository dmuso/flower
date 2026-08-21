package app_test

import (
	"net/http"
	"strings"
	"testing"

	"flower/api/internal/app"
	"flower/api/internal/platform/config"
	"flower/api/internal/testkit"
)

func TestNewRejectsNilConfig(t *testing.T) {
	_, err := app.New(nil)
	if err == nil {
		t.Fatal("expected nil config to fail")
	}
	if !strings.Contains(err.Error(), "config is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLoggerRejectsInvalidLevel(t *testing.T) {
	_, err := app.NewLogger(&config.Config{LogLevel: "loud", Environment: "test"})
	if err == nil {
		t.Fatal("expected invalid log level to fail")
	}
	if !strings.Contains(err.Error(), "invalid LOG_LEVEL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLoggerAcceptsInfoLevel(t *testing.T) {
	logger, err := app.NewLogger(&config.Config{LogLevel: "info", Environment: "test"})
	if err != nil {
		t.Fatalf("newLogger: %v", err)
	}
	if logger == nil {
		t.Fatal("expected logger")
	}
}

func TestStartRejectsNilApp(t *testing.T) {
	var application *app.App
	err := application.Start()
	if err == nil {
		t.Fatal("expected nil app to fail Start")
	}
}

func TestCloseRejectsNilApp(t *testing.T) {
	var application *app.App
	err := application.Close()
	if err == nil {
		t.Fatal("expected nil app to fail Close")
	}
}

func TestCancelOnNilAppDoesNotPanic(t *testing.T) {
	var application *app.App
	application.Cancel()
}

func TestUnverifiedPasswordSignupCannotCreateOrganisation(t *testing.T) {
	h := testkit.New(t)
	res := h.SignUp("maya@example.com", "secret12")
	if res.Code != http.StatusCreated {
		t.Fatalf("signup: %d %s", res.Code, res.Body)
	}
	if h.Count("email_outbox") != 1 {
		t.Fatalf("expected outbox row, got %d", h.Count("email_outbox"))
	}
	login := h.Do(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "maya@example.com", "password": "secret12",
	})
	if login.Code != http.StatusOK {
		t.Fatalf("login unverified: %d %s", login.Code, login.Body)
	}
	cookie := h.SessionCookie(login)
	create := h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie)
	if create.Code != http.StatusForbidden {
		t.Fatalf("org create: %d %s", create.Code, create.Body)
	}
	if testkit.ErrorCode(create) != "email_unverified" {
		t.Fatalf("code: %s body %s", testkit.ErrorCode(create), create.Body)
	}
	if h.Count("organisations") != 0 {
		t.Fatalf("expected no organisations, got %d", h.Count("organisations"))
	}
}

func TestVerifiedPasswordSignupCreatesOrgAndProjectAsOwner(t *testing.T) {
	h := testkit.New(t)
	if res := h.SignUp("maya@example.com", "secret12"); res.Code != http.StatusCreated {
		t.Fatalf("signup: %d %s", res.Code, res.Body)
	}
	cookie := h.Verify("maya@example.com")
	me := h.Do(http.MethodGet, "/api/v1/me", nil, cookie)
	if me.Code != http.StatusOK {
		t.Fatalf("me: %d %s", me.Code, me.Body)
	}
	if testkit.StringField(me.JSON(), "username") != "maya" {
		t.Fatalf("username: %v", me.JSON())
	}
	if me.JSON()["email_verified_at"] == nil {
		t.Fatal("expected email_verified_at")
	}

	orgRes := h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie)
	if orgRes.Code != http.StatusCreated {
		t.Fatalf("org: %d %s", orgRes.Code, orgRes.Body)
	}
	orgID := testkit.StringField(orgRes.JSON(), "id")
	if orgID == "" {
		t.Fatal("missing org id")
	}

	projRes := h.Do(http.MethodPost, "/api/v1/organisations/"+orgID+"/projects", map[string]string{"name": "Trail"}, cookie)
	if projRes.Code != http.StatusCreated {
		t.Fatalf("project: %d %s", projRes.Code, projRes.Body)
	}
	body := projRes.JSON()
	if testkit.StringField(body, "name") != "Trail" {
		t.Fatalf("name: %v", body)
	}
	if testkit.StringField(body, "point_scale") != "linear" {
		t.Fatalf("point_scale: %v", body)
	}
	if body["iteration_length_days"] != float64(7) {
		t.Fatalf("iteration_length_days: %v", body["iteration_length_days"])
	}
	if body["initial_velocity"] != float64(10) {
		t.Fatalf("initial_velocity: %v", body["initial_velocity"])
	}
	if testkit.StringField(body, "timezone") != "Australia/Melbourne" {
		t.Fatalf("timezone: %v", body)
	}
	projectID := testkit.StringField(body, "id")

	var orgRole, projectRole string
	if err := h.DB.QueryRow(`SELECT role FROM organisation_memberships WHERE organisation_id = $1`, orgID).Scan(&orgRole); err != nil {
		t.Fatalf("org role: %v", err)
	}
	if orgRole != "owner" {
		t.Fatalf("org role %q", orgRole)
	}
	if err := h.DB.QueryRow(`SELECT role FROM project_memberships WHERE project_id = $1`, projectID).Scan(&projectRole); err != nil {
		t.Fatalf("project role: %v", err)
	}
	if projectRole != "owner" {
		t.Fatalf("project role %q", projectRole)
	}

	var leftoverTable bool
	if err := h.DB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'iterations'
	)`).Scan(&leftoverTable); err != nil {
		t.Fatalf("iterations table: %v", err)
	}
	if leftoverTable {
		t.Fatal("iterations table must be dropped")
	}
	var leftoverWeeks, leftoverIterationID bool
	if err := h.DB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = 'projects' AND column_name = 'iteration_length_weeks'
	)`).Scan(&leftoverWeeks); err != nil {
		t.Fatalf("weeks column: %v", err)
	}
	if leftoverWeeks {
		t.Fatal("projects.iteration_length_weeks must be dropped")
	}
	if err := h.DB.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = 'stories' AND column_name = 'iteration_id'
	)`).Scan(&leftoverIterationID); err != nil {
		t.Fatalf("iteration_id column: %v", err)
	}
	if leftoverIterationID {
		t.Fatal("stories.iteration_id must be dropped")
	}

	stories := h.Do(http.MethodGet, "/api/v1/projects/"+projectID+"/stories", nil, cookie)
	if stories.Code != http.StatusOK {
		t.Fatalf("stories: %d %s", stories.Code, stories.Body)
	}
	list := stories.JSON()
	rawStories, _ := list["stories"].([]any)
	if len(rawStories) != 0 {
		t.Fatalf("expected empty stories, got %v", rawStories)
	}
	pack, _ := list["pack"].(map[string]any)
	if pack["current_points"] != float64(0) || pack["denominator"] != float64(10) {
		t.Fatalf("pack: %v", pack)
	}
	if pack["velocity_source"] != "initial" {
		t.Fatalf("velocity_source: %v", pack["velocity_source"])
	}
	if pack["current_window_ends_at"] == nil || pack["current_window_ends_at"] == "" {
		t.Fatal("expected computed window end")
	}

	again := h.Do(http.MethodGet, "/api/v1/projects/"+projectID, nil, cookie)
	if again.Code != http.StatusOK {
		t.Fatalf("reload project: %d", again.Code)
	}
	if testkit.StringField(again.JSON(), "id") != projectID {
		t.Fatalf("reload id changed: %v", again.JSON())
	}
	if testkit.StringField(again.JSON(), "organisation_id") != orgID {
		t.Fatalf("reload org changed")
	}
}

func TestMagicLinkNewEmailCreatesVerifiedUser(t *testing.T) {
	h := testkit.New(t)
	req := h.Do(http.MethodPost, "/api/v1/auth/magic-link", map[string]string{"email": "luis@example.com"})
	if req.Code != http.StatusOK {
		t.Fatalf("magic-link: %d %s", req.Code, req.Body)
	}
	if h.Count("email_outbox") != 1 {
		t.Fatalf("outbox %d", h.Count("email_outbox"))
	}
	if h.Count("users") != 0 {
		t.Fatal("user must be created on consume, not request")
	}
	token := h.OutboxToken("luis@example.com", "magic_link")
	consume := h.Do(http.MethodPost, "/api/v1/auth/magic-link/consume", map[string]string{"token": token})
	if consume.Code != http.StatusOK {
		t.Fatalf("consume: %d %s", consume.Code, consume.Body)
	}
	if consume.JSON()["email_verified_at"] == nil {
		t.Fatal("magic link must verify email")
	}
	cookie := h.SessionCookie(consume)
	orgRes := h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie)
	if orgRes.Code != http.StatusCreated {
		t.Fatalf("org: %d %s", orgRes.Code, orgRes.Body)
	}
	orgID := testkit.StringField(orgRes.JSON(), "id")
	proj := h.Do(http.MethodPost, "/api/v1/organisations/"+orgID+"/projects", map[string]string{"name": "Trail"}, cookie)
	if proj.Code != http.StatusCreated {
		t.Fatalf("project: %d %s", proj.Code, proj.Body)
	}
}

func TestPasswordAndMagicLinkLoginLandOnLastProject(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	cookie := h.Verify("maya@example.com")
	orgID := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie).JSON(), "id")
	projectID := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations/"+orgID+"/projects", map[string]string{"name": "Trail"}, cookie).JSON(), "id")
	if h.Do(http.MethodGet, "/api/v1/projects/"+projectID+"/stories", nil, cookie).Code != http.StatusOK {
		t.Fatal("board load")
	}

	h.Do(http.MethodPost, "/api/v1/auth/logout", nil, cookie)
	login := h.Do(http.MethodPost, "/api/v1/auth/login", map[string]string{"email": "maya@example.com", "password": "secret12"})
	if login.Code != http.StatusOK {
		t.Fatalf("login: %d %s", login.Code, login.Body)
	}
	last, _ := login.JSON()["last_project"].(map[string]any)
	if last == nil || last["id"] != projectID {
		t.Fatalf("password login last_project: %v", login.JSON()["last_project"])
	}

	h.Do(http.MethodPost, "/api/v1/auth/logout", nil, h.SessionCookie(login))
	h.Do(http.MethodPost, "/api/v1/auth/magic-link", map[string]string{"email": "maya@example.com"})
	token := h.OutboxToken("maya@example.com", "magic_link")
	magic := h.Do(http.MethodPost, "/api/v1/auth/magic-link/consume", map[string]string{"token": token})
	if magic.Code != http.StatusOK {
		t.Fatalf("magic login: %d %s", magic.Code, magic.Body)
	}
	last, _ = magic.JSON()["last_project"].(map[string]any)
	if last == nil || last["id"] != projectID {
		t.Fatalf("magic login last_project: %v", magic.JSON()["last_project"])
	}
}

func TestCrossTenantProjectIsNotFound(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("a@example.com", "secret12")
	a := h.Verify("a@example.com")
	orgA := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, a).JSON(), "id")
	projectA := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations/"+orgA+"/projects", map[string]string{"name": "Trail"}, a).JSON(), "id")

	h.SignUp("b@example.com", "secret12")
	b := h.Verify("b@example.com")
	orgB := testkit.StringField(h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Other"}, b).JSON(), "id")
	if orgB == "" {
		t.Fatal("org B")
	}

	got := h.Do(http.MethodGet, "/api/v1/projects/"+projectA, nil, b)
	if got.Code != http.StatusNotFound {
		t.Fatalf("cross-tenant project: %d %s", got.Code, got.Body)
	}
	if testkit.ErrorCode(got) != "not_found" {
		t.Fatalf("expected not_found, got %s", testkit.ErrorCode(got))
	}
	stories := h.Do(http.MethodGet, "/api/v1/projects/"+projectA+"/stories", nil, b)
	if stories.Code != http.StatusNotFound {
		t.Fatalf("cross-tenant stories: %d %s", stories.Code, stories.Body)
	}
}

func TestSignupTakenAndBlankValidation(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	taken := h.SignUp("maya@example.com", "secret12")
	if taken.Code != http.StatusConflict {
		t.Fatalf("taken: %d %s", taken.Code, taken.Body)
	}
	blank := h.SignUp("", "")
	if blank.Code != http.StatusBadRequest {
		t.Fatalf("blank: %d %s", blank.Code, blank.Body)
	}
}

func TestPasswordLoginOnMagicLinkOnlyUser(t *testing.T) {
	h := testkit.New(t)
	h.Do(http.MethodPost, "/api/v1/auth/magic-link", map[string]string{"email": "link@example.com"})
	token := h.OutboxToken("link@example.com", "magic_link")
	consume := h.Do(http.MethodPost, "/api/v1/auth/magic-link/consume", map[string]string{"token": token})
	if consume.Code != http.StatusOK {
		t.Fatalf("consume: %d %s", consume.Code, consume.Body)
	}
	login := h.Do(http.MethodPost, "/api/v1/auth/login", map[string]string{"email": "link@example.com", "password": "nope"})
	if login.Code != http.StatusUnauthorized {
		t.Fatalf("login: %d %s", login.Code, login.Body)
	}
	if testkit.ErrorCode(login) != "unauthorized" {
		t.Fatalf("code %s", testkit.ErrorCode(login))
	}
	testkit.MustContain(t, string(login.Body), "magic link")
}

func TestUsernameCollisionGetsSuffix(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("Maya@example.com", "secret12")
	h.SignUp("maya@other.com", "secret12")
	var first, second string
	if err := h.DB.QueryRow(`SELECT username FROM users WHERE email = 'maya@example.com'`).Scan(&first); err != nil {
		t.Fatal(err)
	}
	if err := h.DB.QueryRow(`SELECT username FROM users WHERE email = 'maya@other.com'`).Scan(&second); err != nil {
		t.Fatal(err)
	}
	if first != "maya" {
		t.Fatalf("first username %q", first)
	}
	if second == "maya" || len(second) != len("maya-xxxx") {
		t.Fatalf("second username %q", second)
	}
}

func TestCORSUsesExplicitOriginAndCredentials(t *testing.T) {
	h := testkit.New(t)
	res := h.Do(http.MethodOptions, "/api/v1/me", nil)
	if res.Header.Get("Access-Control-Allow-Origin") != h.Origin {
		t.Fatalf("origin %q", res.Header.Get("Access-Control-Allow-Origin"))
	}
	if res.Header.Get("Access-Control-Allow-Credentials") != "true" {
		t.Fatal("credentials")
	}
	if res.Header.Get("Access-Control-Allow-Origin") == "*" {
		t.Fatal("must not use *")
	}
}

func TestBlankOrganisationAndProjectNames(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	cookie := h.Verify("maya@example.com")
	org := h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "  "}, cookie)
	if org.Code != http.StatusBadRequest {
		t.Fatalf("blank org: %d %s", org.Code, org.Body)
	}
	orgOK := h.Do(http.MethodPost, "/api/v1/organisations", map[string]string{"name": "Acme"}, cookie)
	orgID := testkit.StringField(orgOK.JSON(), "id")
	proj := h.Do(http.MethodPost, "/api/v1/organisations/"+orgID+"/projects", map[string]string{"name": ""}, cookie)
	if proj.Code != http.StatusBadRequest {
		t.Fatalf("blank project: %d %s", proj.Code, proj.Body)
	}
}

func TestResendVerifyEmailWritesVerifyEmailNotMagicLink(t *testing.T) {
	h := testkit.New(t)
	if res := h.SignUp("maya@example.com", "secret12"); res.Code != http.StatusCreated {
		t.Fatalf("signup: %d %s", res.Code, res.Body)
	}
	login := h.Do(http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email": "maya@example.com", "password": "secret12",
	})
	if login.Code != http.StatusOK {
		t.Fatalf("login: %d %s", login.Code, login.Body)
	}
	cookie := h.SessionCookie(login)

	anon := h.Do(http.MethodPost, "/api/v1/auth/verify-email", nil)
	if anon.Code != http.StatusUnauthorized {
		t.Fatalf("anon resend: %d %s", anon.Code, anon.Body)
	}

	res := h.Do(http.MethodPost, "/api/v1/auth/verify-email", nil, cookie)
	if res.Code != http.StatusOK {
		t.Fatalf("resend: %d %s", res.Code, res.Body)
	}

	query := "SELECT kind FROM email_outbox WHERE to_email = $1 ORDER BY created_at ASC"
	rows, err := h.DB.Query(query, "maya@example.com")
	if err != nil {
		t.Fatalf("outbox: %v", err)
	}
	defer rows.Close()
	var kinds []string
	for rows.Next() {
		var kind string
		if err := rows.Scan(&kind); err != nil {
			t.Fatalf("scan: %v", err)
		}
		kinds = append(kinds, kind)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows: %v", err)
	}
	if len(kinds) != 2 {
		t.Fatalf("expected 2 outbox rows, got %v", kinds)
	}
	for _, kind := range kinds {
		if kind != "verify_email" {
			t.Fatalf("expected verify_email, got %q in %v", kind, kinds)
		}
	}
}

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
