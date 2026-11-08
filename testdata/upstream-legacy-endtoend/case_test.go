package endtoend

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type Testcase struct {
	Name       string
	Path       string
	ConfigName string
	Stderr     []byte
}

func parseStderr(t *testing.T, dir, testctx string) []byte {
	t.Helper()
	paths := []string{
		filepath.Join(dir, "stderr", fmt.Sprintf("%s.txt", testctx)),
		filepath.Join(dir, fmt.Sprintf("stderr_%s.txt", runtime.GOOS)),
		filepath.Join(dir, "stderr.txt"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			blob, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			return blob
		}
	}
	return nil
}

func FindTests(t *testing.T, root, testctx string) []*Testcase {
	var tcs []*Testcase
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "sqlc.json" || info.Name() == "sqlc.yaml" || info.Name() == "sqlc.yml" {
			dir := filepath.Dir(path)
			tcs = append(tcs, &Testcase{
				Path:       dir,
				Name:       strings.TrimPrefix(dir, root+string(filepath.Separator)),
				ConfigName: info.Name(),
				Stderr:     parseStderr(t, dir, testctx),
			})
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return tcs
}
