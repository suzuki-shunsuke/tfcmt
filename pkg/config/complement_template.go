package config

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type ComplementTemplateEntry struct {
	Value string
}

func (entry *ComplementTemplateEntry) Type() string {
	return "template"
}

func (entry *ComplementTemplateEntry) Entry() (string, error) {
	tmpl, err := template.New("_").Funcs(sprig.TxtFuncMap()).Parse(entry.Value)
	if err != nil {
		return "", fmt.Errorf("parse a template %s: %w", entry.Value, err)
	}
	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, nil); err != nil {
		return "", fmt.Errorf("render a template with params %s: %w", entry.Value, err)
	}
	return buf.String(), nil
}

func newComplementTemplateEntry(m map[string]interface{}, entry *ComplementTemplateEntry) error {
	v, ok := m["value"]
	if !ok {
		return errors.New(`"value" is required`)
	}
	val, ok := v.(string)
	if !ok {
		return errors.New(`"value" must be string`)
	}
	entry.Value = val
	return nil
}
