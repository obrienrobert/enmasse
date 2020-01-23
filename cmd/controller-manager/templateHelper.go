/*
 * Copyright 2018-2019, EnMasse authors.
 * License: Apache License 2.0 (see the file LICENSE or http://apache.org/licenses/LICENSE-2.0.html).
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"text/template"
)

type TemplateHelper struct {
	TemplatePath string
	TemplateList []string
}

// Creates a new templates helper and populates the values for all
// templates properties. Some of them (like the hostname) are set
// by the user in the custom resource
func NewTemplateHelper() *TemplateHelper {

	log.Info("Installing Monitoring Resources")

	templatePath := "./templates/" // Deployments
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		templatePath = "../templates/prometheus-rules" // Local
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			panic("cannot find templates")
		}
	}

	return &TemplateHelper{
		TemplatePath: templatePath,
		TemplateList: GetTemplateList(),
	}
}

func GetTemplateList() []string {
	templateList := []string{
		"010-PrometheusRules-kube-metrics.yaml",
	}
	return templateList
}

// load a templates from a given resource name. The templates must be located
// under ./templates and the filename must be <resource-name>.yaml
func (h *TemplateHelper) loadTemplate(name string) ([]byte, error) {
	path := fmt.Sprintf("%s/%s", h.TemplatePath, name)
	tpl, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parser := template.New(path)

	parsed, err := parser.Parse(string(tpl))
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = parsed.Execute(&buffer, nil)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (h *TemplateHelper) CreateResource(template string) (runtime.Object, error) {
	tpl, err := h.loadTemplate(template)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resource := unstructured.Unstructured{}
	err = yaml.Unmarshal(tpl, &resource)

	if err != nil {
		return nil, err
	}

	return &resource, nil
}