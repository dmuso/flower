package story

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"flower/api/internal/ports"
	"flower/api/internal/types"
)

type denyAccess struct{}

func (denyAccess) OrganisationRole(context.Context, string, string) (string, error) {
	return "", types.ErrNotFound
}

func (denyAccess) OpenProject(context.Context, string, string) (string, string, error) {
	return "", "", types.ErrNotFound
}

type unusedPlanner struct{}

func (unusedPlanner) ForProject(context.Context, string) (ports.Pack, error) {
	return ports.Pack{}, nil
}

func TestListUnknownProjectIsNotFound(t *testing.T) {
	svc := NewService(nil, denyAccess{}, unusedPlanner{})
	_, err := svc.List(context.Background(), "user-1", "project-1")
	if !errors.Is(err, types.ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestStorySourceDoesNotImportPlanningOrProjects(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("caller")
	}
	entries, err := os.ReadDir(filepath.Dir(file))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(filepath.Dir(file), e.Name()))
		if err != nil {
			t.Fatal(err)
		}
		src := string(b)
		if strings.Contains(src, "domain/planning") {
			t.Fatalf("%s imports domain/planning", e.Name())
		}
		if strings.Contains(src, "FROM projects") {
			t.Fatalf("%s SELECTs FROM projects", e.Name())
		}
	}
}
