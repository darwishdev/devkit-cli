package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type TemplatesInterface interface {
	LoadTemplate(path string, data interface{}) (bytes.Buffer, error)
	LoadLayerTemplates(pattern string, data interface{}) (map[string]bytes.Buffer, error)
}

type Templates struct {
	basePath string
}

func NewTemplates(basePath string) TemplatesInterface {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPathDir := filepath.Base(exePath)
	path := fmt.Sprintf("%s/%s", execPathDir, basePath)
	return &Templates{
		basePath: path,
	}
}
func (t *Templates) LoadTemplate(path string, data interface{}) (bytes.Buffer, error) {
	tmplBytes, err := os.ReadFile(path)
	var tplBuffer bytes.Buffer
	if err != nil {
		return tplBuffer, err
	}
	tmpl, err := template.New("tmpl").Parse(string(tmplBytes))
	if err != nil {
		return tplBuffer, err
	}
	err = tmpl.Execute(&tplBuffer, data)
	if err != nil {
		return tplBuffer, err
	}
	return tplBuffer, nil
}

func (t *Templates) LoadLayerTemplates(pattern string, data interface{}) (map[string]bytes.Buffer, error) {
	result := make(map[string]bytes.Buffer)
	matches, err := filepath.Glob(fmt.Sprintf("%s%s", t.basePath, pattern))
	if err != nil {
		return nil, err
	}
	for _, templatePath := range matches {
		templateReader, err := t.LoadTemplate(templatePath, data)
		if err != nil {
			return nil, err
		}
		fileNameWithExt := filepath.Base(templatePath)
		fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))
		nameParts := strings.Split(fileName, "_")
		resultKey := nameParts[len(nameParts)-1]
		result[resultKey] = templateReader
	}
	return result, nil
}
