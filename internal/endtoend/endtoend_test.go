package endtoend

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

var update = flag.Bool("update", false, "update expected output files")

// sqlcConfig represents the relevant parts of sqlc.yaml/json
type sqlcConfig struct {
	SQL []sqlBlock `json:"sql" yaml:"sql"`
}

type sqlBlock struct {
	Codegen []codegenBlock `json:"codegen" yaml:"codegen"`
}

type codegenBlock struct {
	Out string `json:"out" yaml:"out"`
}

func TestEndToEnd(t *testing.T) {
	root := filepath.Join("testdata")

	testcases := FindTests(t, root, "base")
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			// sqlc caches the shared local WASM plugin during generation.
			runTest(t, tc)
		})
	}
}

func runTest(t *testing.T, tc *Testcase) {
	t.Helper()

	absPath, err := filepath.Abs(tc.Path)
	if err != nil {
		t.Fatalf("failed to get absolute path: %v", err)
	}

	// Parse the config to find output directories
	configPath := filepath.Join(absPath, tc.ConfigName)
	outDirs, err := parseOutputDirs(configPath)
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	if len(outDirs) == 0 {
		t.Fatal("no output directories found in config")
	}

	// Read expected files from all output directories before generation (skip if updating)
	expectedFiles := make(map[string]map[string][]byte)
	if !*update {
		for _, outDir := range outDirs {
			fullOutDir := filepath.Join(absPath, outDir)
			if _, err := os.Stat(fullOutDir); os.IsNotExist(err) {
				continue
			}
			files, err := readDirFiles(fullOutDir)
			if err != nil {
				t.Fatalf("failed to read expected files from %s: %v", outDir, err)
			}
			expectedFiles[outDir] = files
		}
	}

	// Run sqlc generate
	cmd := exec.Command("go", "tool", "sqlc", "generate")
	cmd.Dir = absPath

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	// Check if we expected an error
	if tc.Stderr != nil {
		if !bytes.Contains(stderr.Bytes(), bytes.TrimSpace(tc.Stderr)) {
			t.Errorf("expected stderr to contain:\n%s\ngot:\n%s", tc.Stderr, stderr.String())
		}
		return
	}

	if err != nil {
		t.Fatalf("sqlc generate failed: %v\nstdout: %s\nstderr: %s", err, stdout.String(), stderr.String())
	}

	// If updating, we're done - the generated files are now the expected output
	if *update {
		t.Logf("updated expected output in %v", outDirs)
		return
	}

	// Compare generated files with expected for each output directory
	for outDir, expected := range expectedFiles {
		fullOutDir := filepath.Join(absPath, outDir)
		generated, err := readDirFiles(fullOutDir)
		if err != nil {
			t.Fatalf("failed to read generated files from %s: %v", outDir, err)
		}

		// Compare expected vs generated
		for name, expectedContent := range expected {
			generatedContent, ok := generated[name]
			if !ok {
				t.Errorf("[%s] expected file %s was not generated", outDir, name)
				continue
			}

			if !bytes.Equal(expectedContent, generatedContent) {
				t.Errorf("[%s] file %s differs from expected:\n--- expected\n%s\n--- got\n%s",
					outDir, name, string(expectedContent), string(generatedContent))
			}
		}

		// Check for unexpected files
		for name := range generated {
			if _, ok := expected[name]; !ok {
				t.Errorf("[%s] unexpected file generated: %s", outDir, name)
			}
		}
	}
}

// parseOutputDirs reads the sqlc config and returns the output directories
func parseOutputDirs(configPath string) ([]string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config sqlcConfig

	ext := filepath.Ext(configPath)
	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	}

	var outDirs []string
	for _, sql := range config.SQL {
		for _, cg := range sql.Codegen {
			if cg.Out != "" {
				outDirs = append(outDirs, cg.Out)
			}
		}
	}

	return outDirs, nil
}

// readDirFiles reads all files in a directory and returns a map of filename to content
func readDirFiles(dir string) (map[string][]byte, error) {
	files := make(map[string][]byte)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		files[relPath] = content
		return nil
	})

	return files, err
}
