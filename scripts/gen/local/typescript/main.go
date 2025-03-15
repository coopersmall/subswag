package main

import (
	"context"
	"fmt"
	"github.com/coopersmall/subswag/utils"
	"github.com/gzuidhof/tygo/tygo"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	ctx := context.Background()
	logger := utils.GetLogger("generate_typescript")
	absPath, err := os.Getwd()
	if err != nil {
		logger.Error(ctx, "Failed to get working directory", err, nil)
		os.Exit(1)
	}

	domainDir := filepath.Join(absPath, "domain")
	logger.Info(ctx, "Domain directory", map[string]interface{}{"path": domainDir})

	var packageConfigs []*tygo.PackageConfig
	err = filepath.Walk(domainDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			relPath, err := filepath.Rel(domainDir, path)
			if err != nil {
				return err
			}

			// if the dir contains the file .tsignore, skip it

			files, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			for _, file := range files {
				if file.Name() == ".tsignore" {
					return filepath.SkipDir
				}
			}

			// Get the directory name in lowercase
			dirName := strings.ToLower(filepath.Base(relPath))

			packagePath := filepath.Join("github.com/coopersmall/subswag/domain", relPath)
			// Create output path with lowercase directory name
			outputPath := filepath.Join("./frontend/src/types", dirName+".generated.ts")

			if dirName == "." {
				packagePath = "github.com/coopersmall/subswag/domain"
				outputPath = "./frontend/src/types/domain.generated.ts"
			}

			excludedFiles, err := getExcludedFiles(path)
			if err != nil {
				return err
			}
			packageConfigs = append(packageConfigs, &tygo.PackageConfig{
				Path:       packagePath,
				OutputPath: outputPath,
				TypeMappings: map[string]string{
					"time.Time":              "string /* ISO8601 */",
					"*time.Time":             "string | null /* ISO8601 */",
					"map[string]interface{}": "Record<string, unknown>",
					"utils.ID":               "number",
					"utils.Time":             "string /* ISO8601 */",
					"*utils.Time":            "string | null /* ISO8601 */",
				},
				FallbackType: "unknown",
				Flavor:       "yaml",
				ExcludeFiles: excludedFiles,
			})
		}
		return nil
	})
	if err != nil {
		logger.Error(ctx, "Failed to walk domain directory", err, nil)
		os.Exit(1)
	}

	config := &tygo.Config{
		Packages: packageConfigs,
	}

	generator := tygo.New(config)
	err = generator.Generate()
	if err != nil {
		logger.Error(ctx, "Failed to generate TypeScript definitions", err, nil)
		os.Exit(1)
	}

	logger.Info(ctx, "TypeScript definitions generated successfully", map[string]any{
		"packageConfigs": packageConfigs,
	})
}

func getExcludedFiles(dir string) ([]string, error) {
	var excludedFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dir {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		re := regexp.MustCompile(`//ts:ignore`)
		if re.Match(content) {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			excludedFiles = append(excludedFiles, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}
	return excludedFiles, nil
}
