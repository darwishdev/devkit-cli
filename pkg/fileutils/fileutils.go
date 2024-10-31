package fileutils

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FileUtilsInterface interface {
	ReadExcelFile(filePath string) (*excelize.File, *bytes.Buffer, error)
	ReplaceFile(filePath string, oldText string, newText string) error
	AppendToFile(filePath string, templateContent bytes.Buffer) error
	CopyFiles(globPattern string, mapper func(string) string) error
	CopyFile(src string, dst string) error
	ReplaceAll(filePath string, oldText, newText string) error
	ReplaceMultiple(filePath string, replacements map[string]string) error
	ReplaceInDir(dirPath string, replacements map[string]string) error
}

type FileUtils struct {
}

func NewFileUtils() FileUtilsInterface {
	return &FileUtils{}
}
func (f *FileUtils) ReadExcelFile(filePath string) (*excelize.File, *bytes.Buffer, error) {
	excelFile, err := os.ReadFile(filePath) // Replace with your Excel file path
	if err != nil {
		return nil, nil, err
	}
	fileBuffer := bytes.NewBuffer(excelFile)
	excelBufer := bytes.NewBuffer(excelFile)
	file, err := excelize.OpenReader(excelBufer)
	if err != nil {
		return nil, fileBuffer, err
	}
	return file, fileBuffer, nil
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
func (f *FileUtils) CopyFile(src string, dst string) error {
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
		if d.Name() == ".git" {
			return filepath.SkipDir
		}
		if !d.IsDir() {
			err := f.ReplaceMultiple(path, replacements)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
