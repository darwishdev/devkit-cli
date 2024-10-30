package templates

import (
	"bytes"
	"os"
	"text/template"
)

type TemplatesInterface interface {
	LoadTemplate(path string, data interface{}) (bytes.Buffer, error)
}

type Templates struct {
}

func NewTemplates() TemplatesInterface {
	return &Templates{}
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
