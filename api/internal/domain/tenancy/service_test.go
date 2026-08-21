package tenancy_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"flower/api/internal/testkit"
	"flower/api/internal/types"
)

func TestOrganisationRoleUnknownIsNotFound(t *testing.T) {
	h := testkit.New(t)
	h.SignUp("maya@example.com", "secret12")
	cookie := h.Verify("maya@example.com")
	missing := h.Do(http.MethodPost, "/api/v1/organisations/00000000-0000-0000-0000-000000000001/projects", map[string]string{"name": "Trail"}, cookie)
	if missing.Code != http.StatusNotFound {
		t.Fatalf("status %d %s", missing.Code, missing.Body)
	}
	if !errors.Is(types.ErrNotFound, types.ErrNotFound) {
		t.Fatal("sanity")
	}
	_ = context.Background()
}
