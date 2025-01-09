package component

import (
	"bytes"
	"html/template"
)

type Component struct {
	Template string
}

func (c *Component) Build(data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(c.Template)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
