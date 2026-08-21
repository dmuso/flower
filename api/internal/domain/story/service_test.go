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
	dir := filepath.Dir(file)
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	planningImport := strings.Join([]string{"domain", "planning"}, "/")
	projectsSQL := strings.Join([]string{"FROM", "projects"}, " ")
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".go") {
			continue
		}
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			t.Fatal(err)
		}
		src := string(b)
		if strings.Contains(src, planningImport) {
			t.Fatalf("%s imports planning", name)
		}
		if strings.Contains(src, projectsSQL) {
			t.Fatalf("%s SELECTs projects", name)
		}
	}
}
