package actions

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const SCHEMA_FILE = "db/sql/schema.sql"

func GetSchema() (string, error) {
	// Get the current file's directory
	_, currentFile, _, _ := runtime.Caller(0)

	// Navigate up until finding go.mod
	dir := filepath.Dir(currentFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find go.mod file")
		}
		dir = parent
	}

	// Open and read the file
	file, err := os.Open(filepath.Join(dir, SCHEMA_FILE))
	if err != nil {
		return "", err
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

