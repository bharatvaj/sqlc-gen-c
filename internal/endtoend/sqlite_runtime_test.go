package endtoend

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// Exercise the checked-in sqlc-generated fixture against a real SQLite connection.
func TestSQLiteDirectRuntime(t *testing.T) {
	dir := filepath.Join("testdata", "direct", "sqlite")
	binary := filepath.Join(t.TempDir(), "sqlite-test")
	cmd := exec.Command("cc", "-std=c99", "-Wall", "-Wextra", "-Werror", "runtime.c", "gen/queries.c", "-lsqlite3", "-o", binary)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("compile SQLite fixture: %v\n%s", err, out)
	}
	if out, err := exec.Command(binary).CombinedOutput(); err != nil {
		t.Fatalf("run SQLite fixture: %v\n%s", err, out)
	}
}
