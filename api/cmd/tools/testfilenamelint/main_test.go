package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLintTestFileNames(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		want  []string
	}{
		{
			name:  "pairing allowed",
			files: []string{"handler.go", "handler_test.go"},
		},
		{
			name:  "pairing allowed in subdirectory",
			files: []string{"internal/app/app.go", "internal/app/app_test.go"},
		},
		{
			name:  "integration allowed without sibling",
			files: []string{"auth_integration_test.go"},
		},
		{
			name:  "perf allowed without sibling",
			files: []string{"rank_perf_test.go"},
		},
		{
			name:  "extra suffix rejected",
			files: []string{"slice0_test.go", "slice0_more_test.go", "handler_helpers_test.go"},
			want:  []string{"handler_helpers_test.go", "slice0_more_test.go", "slice0_test.go"},
		},
		{
			name:  "skips git and tmp",
			files: []string{".git/unpaired_test.go", "tmp/unpaired_test.go", "ok.go", "ok_test.go"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := writeTree(t, tc.files)
			got, err := lintTestFileNames(root)
			if err != nil {
				t.Fatalf("lintTestFileNames: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("issues = %#v, want %#v", got, tc.want)
			}
		})
	}
}

func TestWriteIssues(t *testing.T) {
	var buf bytes.Buffer
	if err := writeIssues(&buf, []string{"handler_helpers_test.go"}); err != nil {
		t.Fatal(err)
	}
	want := "handler_helpers_test.go: unit test files must be named <code_file>_test.go with a same-directory <code_file>.go, or use <domain>_integration_test.go / <domain>_perf_test.go.\nTest file name lint failed with 1 issue(s).\n"
	if buf.String() != want {
		t.Fatalf("got %q", buf.String())
	}
}

func writeTree(t *testing.T, files []string) string {
	t.Helper()
	root := t.TempDir()
	for _, rel := range files {
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("package p\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}
