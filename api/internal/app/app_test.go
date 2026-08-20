package app

import (
	"strings"
	"testing"

	"flower/api/internal/platform/config"
)

func TestNewRejectsNilConfig(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected nil config to fail")
	}
	if !strings.Contains(err.Error(), "config is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLoggerRejectsInvalidLevel(t *testing.T) {
	_, err := newLogger(&config.Config{LogLevel: "loud", Environment: "test"})
	if err == nil {
		t.Fatal("expected invalid log level to fail")
	}
	if !strings.Contains(err.Error(), "invalid LOG_LEVEL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewLoggerAcceptsInfoLevel(t *testing.T) {
	logger, err := newLogger(&config.Config{LogLevel: "info", Environment: "test"})
	if err != nil {
		t.Fatalf("newLogger: %v", err)
	}
	if logger == nil {
		t.Fatal("expected logger")
	}
}

func TestStartRejectsNilApp(t *testing.T) {
	var application *App
	err := application.Start()
	if err == nil {
		t.Fatal("expected nil app to fail Start")
	}
}

func TestCloseRejectsNilApp(t *testing.T) {
	var application *App
	err := application.Close()
	if err == nil {
		t.Fatal("expected nil app to fail Close")
	}
}

func TestCancelOnNilAppDoesNotPanic(t *testing.T) {
	var application *App
	application.Cancel()
}
