package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetMigrationsDirPointsAtMigrationFiles(t *testing.T) {
	dir := getMigrationsDir()
	up := filepath.Join(dir, "000001_create_core_schema.up.sql")
	info, err := os.Stat(up)
	if err != nil {
		t.Fatalf("expected migration file at %s: %v", up, err)
	}
	if info.IsDir() {
		t.Fatalf("expected %s to be a file", up)
	}
}
