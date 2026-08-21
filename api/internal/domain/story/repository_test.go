package story

import (
	"os"
	"strings"
	"testing"
)

func TestRepositorySQLDoesNotMentionProjectsOrIterations(t *testing.T) {
	body, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatal(err)
	}
	src := strings.ToLower(string(body))
	if strings.Contains(src, "projects") {
		t.Fatal("story repository SQL must not mention projects")
	}
	if strings.Contains(src, "iterations") {
		t.Fatal("story repository SQL must not mention iterations")
	}
}
