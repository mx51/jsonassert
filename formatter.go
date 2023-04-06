package jsonassert

import (
	"bytes"
	"testing"
	stdtemplate "text/template"
)

func Format(t *testing.T, template string, data any) string {
	templ := stdtemplate.New("")
	if _, err := templ.Parse(template); err != nil {
		t.Error("jsonassert.Format: could not parse template: %w", err)
	}

	var res bytes.Buffer
	if err := templ.Execute(&res, data); err != nil {
		t.Error("jsonassert.Format: could not render template: %w", err)
	}

	return res.String()
}
