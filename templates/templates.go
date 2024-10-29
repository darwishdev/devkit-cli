package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/darwishdev/devkit-cli/config"
	"github.com/darwishdev/devkit-cli/models"
)

type TemplatesInterface interface {
	LoadEndpointTemplates(domainName string, featureName string, endpintName string, data models.ConfigMap) (map[string]bytes.Buffer, error)
	LoadLayerTemplates(prefix string, data models.ConfigMap) (map[string]bytes.Buffer, error)
	LoadDomainTemplates(domainName string) (map[string]bytes.Buffer, error)
	ParseTemplateName(path string) string
	LoadFeatureTemplates(domainName string, featureName string) (map[string]bytes.Buffer, error)
	LoadTemplate(path string, data interface{}) (bytes.Buffer, error)
}

type Templates struct {
	config         config.ConfigInterface
	basePath       string
	domainPrefix   string
	featurePrefix  string
	endpointPrefix string
}

func NewTemplates(config config.ConfigInterface) TemplatesInterface {
	basePath := fmt.Sprintf("%s/templates/tmpls", config.GetConfig().ExePath)
	return &Templates{
		config:         config,
		basePath:       basePath,
		domainPrefix:   "domain",
		featurePrefix:  "feature",
		endpointPrefix: "endpoint",
	}
}

func (t *Templates) ParseTemplateName(path string) string {
	removedExt := strings.Replace(path, ".tmpl", "", 1)
	slashParts := strings.Split(removedExt, "/")
	return slashParts[len(slashParts)-1]
}
func (t *Templates) LoadTemplate(path string, data interface{}) (bytes.Buffer, error) {
	name := t.ParseTemplateName(path)
	tmplBytes, err := os.ReadFile(path)
	var tplBuffer bytes.Buffer
	if err != nil {
		return tplBuffer, err
	}
	tmpl, err := template.New(name).Parse(string(tmplBytes))
	if err != nil {
		return tplBuffer, err
	}
	err = tmpl.Execute(&tplBuffer, data)
	if err != nil {
		return tplBuffer, err
	}
	return tplBuffer, nil
}

func (t *Templates) LoadLayerTemplates(prefix string, data models.ConfigMap) (map[string]bytes.Buffer, error) {
	result := make(map[string]bytes.Buffer)
	pattern := fmt.Sprintf("%s/%s*", t.basePath, prefix)
	matches, err := filepath.Glob(pattern)
	fmt.Println("pattern", pattern, "matches ", matches)
	if err != nil {
		return nil, err
	}
	for _, templatePath := range matches {
		templateReader, err := t.LoadTemplate(templatePath, data)
		if err != nil {
			return nil, err
		}
		name := t.ParseTemplateName(templatePath)
		nameParts := strings.Split(name, "_")
		resultKey := nameParts[len(nameParts)-1]
		result[resultKey] = templateReader
	}
	return result, nil
}
func (t *Templates) LoadDomainTemplates(domainName string) (map[string]bytes.Buffer, error) {
	data := t.config.WithDomain(domainName)
	fmt.Println("domain", data)
	result, err := t.LoadLayerTemplates("domain", data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Templates) LoadFeatureTemplates(domainName string, featureName string) (map[string]bytes.Buffer, error) {
	data := t.config.WithFeature(domainName, featureName)
	result, err := t.LoadLayerTemplates("feature", data)
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (t *Templates) LoadEndpointTemplates(domainName string, featureName string, endpintName string, data models.ConfigMap) (map[string]bytes.Buffer, error) {
	result, err := t.LoadLayerTemplates("endpoint", data)
	if err != nil {
		return nil, err
	}
	return result, nil

}
