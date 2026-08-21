package project

import (
	"os"
	"strings"
	"testing"
)

func TestRepositoryDoesNotWriteIterationLengthWeeks(t *testing.T) {
	body, err := os.ReadFile("repository.go")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "iteration_length_weeks") {
		t.Fatal("project repository must not write leftover iteration_length_weeks")
	}
}
