package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutesRejectsNilRouter(t *testing.T) {
	err := SetupRoutes(nil, &Dependencies{Version: "dev"})
	if err == nil {
		t.Fatal("expected nil router to fail")
	}
}

func TestSetupRoutesRejectsNilDependencies(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := SetupRoutes(gin.New(), nil)
	if err == nil {
		t.Fatal("expected nil dependencies to fail")
	}
}

func TestSetupRoutesRejectsEmptyVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := SetupRoutes(gin.New(), &Dependencies{})
	if err == nil {
		t.Fatal("expected empty version to fail")
	}
}

func TestSetupRoutesRegistersCoreEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	if err := SetupRoutes(router, &Dependencies{Version: "dev", Pinger: stubPinger{}}); err != nil {
		t.Fatalf("SetupRoutes: %v", err)
	}

	for _, path := range []string{"/health", "/ready", "/api/version"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s status: got %d, want %d", path, rec.Code, http.StatusOK)
		}
	}
}
