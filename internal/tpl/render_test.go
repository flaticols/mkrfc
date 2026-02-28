package tpl_test

import (
	"strings"
	"testing"

	"github.com/flaticols/mkrfc/internal/tpl"
)

func TestRender(t *testing.T) {
	tmpl := "Hello, {{ .Name }}! Date: {{ .Date }}"
	data := map[string]string{"Name": "World", "Date": "2024-01-01"}

	out, err := tpl.Render(tmpl, data)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Hello, World!") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "2024-01-01") {
		t.Errorf("expected date in output: %s", out)
	}
}

func TestRenderInvalidTemplate(t *testing.T) {
	_, err := tpl.Render("{{ .Foo", nil)
	if err == nil {
		t.Error("expected error for invalid template")
	}
}
