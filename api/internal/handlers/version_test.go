package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestVersionHandlerReturnsApplicationNameAndVersion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/version", VersionHandler("0.1.0"))

	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["name"] != "flower" {
		t.Fatalf("name: got %q, want %q", body["name"], "flower")
	}
	if body["version"] != "0.1.0" {
		t.Fatalf("version: got %q, want %q", body["version"], "0.1.0")
	}
	if rec.Header().Get("Cache-Control") == "" {
		t.Fatal("expected Cache-Control header")
	}
}
