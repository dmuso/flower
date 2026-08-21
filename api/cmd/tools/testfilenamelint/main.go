package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const issueMessage = "unit test files must be named <code_file>_test.go with a same-directory <code_file>.go, or use <domain>_integration_test.go / <domain>_perf_test.go."

func main() {
	root := "."
	if len(os.Args) > 1 && os.Args[1] != "" {
		root = os.Args[1]
	}

	issues, err := lintTestFileNames(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := writeIssues(os.Stderr, issues); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if len(issues) > 0 {
		os.Exit(1)
	}
}

func lintTestFileNames(root string) ([]string, error) {
	var issues []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "tmp" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.go") {
			return nil
		}
		if allowedTestFileName(path, d.Name()) {
			return nil
		}
		rel := path
		if r, relErr := filepath.Rel(root, path); relErr == nil {
			rel = r
		}
		issues = append(issues, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk test files: %w", err)
	}
	sort.Strings(issues)
	return issues, nil
}

func allowedTestFileName(path, name string) bool {
	if strings.HasSuffix(name, "_integration_test.go") || strings.HasSuffix(name, "_perf_test.go") {
		return true
	}
	stem := strings.TrimSuffix(name, "_test.go")
	sibling := filepath.Join(filepath.Dir(path), stem+".go")
	_, err := os.Stat(sibling)
	return err == nil
}

func writeIssues(w io.Writer, issues []string) error {
	for _, path := range issues {
		if _, err := fmt.Fprintf(w, "%s: %s\n", path, issueMessage); err != nil {
			return err
		}
	}
	if len(issues) > 0 {
		if _, err := fmt.Fprintf(w, "Test file name lint failed with %d issue(s).\n", len(issues)); err != nil {
			return err
		}
	}
	return nil
}
