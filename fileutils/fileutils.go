package fileutils

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/darwishdev/devkit-cli/config"
	"github.com/darwishdev/devkit-cli/models"
)

type FileUtilsInterface interface {
	ReplaceInDirFiles(dirPath string, replacements map[string]string) error
	ReplaceFile(filePath string, oldText string, newText string) error
	AppendToFile(filePath string, templateContent bytes.Buffer) error
	CollectFeatureFiles(domainName string, featureName string) models.ConfigMap
	CollectDomainDirecories(domainName string) models.ConfigMap
	CollectDomainFiles(domainName string) models.ConfigMap
	CopyFile(src, dst string) error
	CopyFiles(globPattern string, mapper func(string) string) error
	ReplaceAll(filePath string, oldText, newText string) error
	ReplaceMultiple(filePath string, replacements map[string]string) error
	ReplaceInDir(dirPath string, replacements map[string]string) error
}

type FileUtils struct {
	config config.ConfigInterface
	AppDir string
}

func NewFileUtils(config config.ConfigInterface) FileUtilsInterface {
	return &FileUtils{
		config: config,
		AppDir: "app",
	}
}
func (f *FileUtils) CollectDomainDirecories(domainName string) models.ConfigMap {
	basePath := fmt.Sprintf("%s/%s", f.AppDir, domainName)
	return models.ConfigMap{
		"adapter": fmt.Sprintf("%s/adapter", basePath),
		"repo":    fmt.Sprintf("%s/repo", basePath),
		"usecase": fmt.Sprintf("%s/usecase", basePath),
	}
}

func (f *FileUtils) CollectDomainFiles(domainName string) models.ConfigMap {
	basePath := fmt.Sprintf("%s/%s", f.AppDir, domainName)
	return models.ConfigMap{
		"adapter": fmt.Sprintf("%s/adapter/adapter.go", basePath),
		"repo":    fmt.Sprintf("%s/repo/repo.go", basePath),
		"usecase": fmt.Sprintf("%s/usecase/usecase.go", basePath),
	}
}

func (f *FileUtils) CollectFeatureFiles(domainName string, featureName string) models.ConfigMap {
	basePath := fmt.Sprintf("%s/%s", f.AppDir, domainName)
	config := f.config.GetConfig()
	return models.ConfigMap{
		"adapter": fmt.Sprintf("%s/adapter/%s_adapter.go", basePath, featureName),
		"repo":    fmt.Sprintf("%s/repo/%s_repo.go", basePath, featureName),
		"usecase": fmt.Sprintf("%s/usecase/%s_usecase.go", basePath, featureName),
		"proto":   fmt.Sprintf("%s/%s_%s.proto", config.ProtoPath, domainName, featureName),
		"query":   fmt.Sprintf("%s/%s_%s.sql", config.QueryPath, domainName, featureName),
		"api":     fmt.Sprintf("api/%s_%s_rpc.go", domainName, featureName),
	}
}
func (f *FileUtils) AppendToFile(filePath string, templateContent bytes.Buffer) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()
	// 2. Append the content
	if _, err := file.Write(templateContent.Bytes()); err != nil {
		fmt.Println("Error appending to file:", err)
		return err
	}
	return nil
}
func (f *FileUtils) ReplaceFile(filePath string, oldText string, newText string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	fileContent := string(file)

	newContent := strings.Replace(fileContent, oldText, newText, 1)
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	return nil
}

// CopyFile copies a single file from src to dst
func (f *FileUtils) CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// CopyFiles copies multiple files matching the given glob pattern.
// The mapper function is used to transform the original filename to the new filename.
func (f *FileUtils) CopyFiles(globPattern string, mapper func(string) string) error {
	files, err := filepath.Glob(globPattern)
	if err != nil {
		return err
	}

	for _, file := range files {
		newFilename := mapper(file)
		err := f.CopyFile(file, newFilename)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReplaceAll replaces all occurrences of oldText with newText in the file.
func (f *FileUtils) ReplaceAll(filePath string, oldText, newText string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := string(file)

	newContent := strings.ReplaceAll(fileContent, oldText, newText)
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReplaceMultiple replaces multiple strings in a file based on a map of old and new strings.
func (f *FileUtils) ReplaceMultiple(filePath string, replacements map[string]string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	fileContent := string(file)

	for oldText, newText := range replacements {
		fileContent = strings.ReplaceAll(fileContent, oldText, newText)
	}

	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReplaceInDir recursively replaces strings in files within a directory based on a map of replacements.
func (f *FileUtils) ReplaceInDir(dirPath string, replacements map[string]string) error {
	return filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error { // Use filepath.WalkDir
		if err != nil {
			return err
		}
		// Ignore .git directory
		if d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() { // Check if it's not a directory
			err := f.ReplaceMultiple(path, replacements)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
func (f *FileUtils) ReplaceInDirFiles(dirPath string, replacements map[string]string) error {
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := f.ReplaceMultiple(path, replacements)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
