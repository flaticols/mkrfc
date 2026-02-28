package tpl

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
)

type Template struct {
	Name    string
	Content string
	Builtin bool
}

func ListTemplates(embeddedFS fs.FS, customDir string) ([]Template, error) {
	data, err := fs.ReadFile(embeddedFS, "templates/default.md")
	if err != nil {
		return nil, err
	}

	templates := []Template{
		{Name: "default", Content: string(data), Builtin: true},
	}

	if _, statErr := os.Stat(customDir); os.IsNotExist(statErr) {
		return templates, nil
	}

	entries, err := filepath.Glob(filepath.Join(customDir, "*.md"))
	if err != nil {
		return nil, err
	}

	for _, path := range entries {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		name := filepath.Base(path)
		name = name[:len(name)-3] // strip .md
		templates = append(templates, Template{Name: name, Content: string(content)})
	}

	return templates, nil
}

func Render(tmplContent string, data any) (string, error) {
	t, err := template.New("rfc").Parse(tmplContent)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
