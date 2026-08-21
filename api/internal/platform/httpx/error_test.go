package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"flower/api/internal/types"

	"github.com/gin-gonic/gin"
)

func TestFromMapsKnownErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		err    error
		status int
		code   string
	}{
		{types.Wrap(types.ErrValidation, "Name the organisation."), http.StatusBadRequest, "validation_failed"},
		{types.ErrEmailTaken, http.StatusConflict, "email_taken"},
		{types.ErrEmailUnverified, http.StatusForbidden, "email_unverified"},
		{types.ErrUseMagicLink, http.StatusUnauthorized, "unauthorized"},
		{types.ErrUnauthorized, http.StatusUnauthorized, "unauthorized"},
		{types.ErrForbidden, http.StatusForbidden, "forbidden"},
		{types.ErrNotFound, http.StatusNotFound, "not_found"},
		{types.ErrTokenExpired, http.StatusUnauthorized, "unauthorized"},
		{types.ErrTokenConsumed, http.StatusUnauthorized, "unauthorized"},
		{types.ErrConflict, http.StatusConflict, "conflict"},
		{errors.New("boom"), http.StatusInternalServerError, "internal_error"},
		{types.ErrValidation, http.StatusBadRequest, "validation_failed"},
	}
	for _, tc := range cases {
		router := gin.New()
		router.GET("/", func(c *gin.Context) { From(c, tc.err) })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != tc.status {
			t.Fatalf("%v status %d want %d", tc.err, rec.Code, tc.status)
		}
		var body map[string]map[string]string
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatal(err)
		}
		if body["error"]["code"] != tc.code {
			t.Fatalf("%v code %s", tc.err, body["error"]["code"])
		}
	}
}

func TestWrite(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/", func(c *gin.Context) { Write(c, 418, "teapot", "short and stout") })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != 418 {
		t.Fatalf("status %d", rec.Code)
	}
}
