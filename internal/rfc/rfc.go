package rfc

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Status string

const (
	StatusDraft      Status = "DRAFT"
	StatusAccepted   Status = "ACCEPTED"
	StatusRejected   Status = "REJECTED"
	StatusSuperseded Status = "SUPERSEDED"
)

type Frontmatter struct {
	Title        string   `yaml:"title"`
	Authors      []string `yaml:"authors"`
	Reviewers    []string `yaml:"reviewers"`
	Status       Status   `yaml:"status"`
	Date         string   `yaml:"date"`
	ResolvedDate string   `yaml:"resolved_date,omitempty"`
	PR           string   `yaml:"pr,omitempty"`
	Tags         []string `yaml:"tags"`
}

// RFC fields ordered to minimize GC-scanned pointer bytes (fieldalignment).
type RFC struct {
	Body     string
	FilePath string
	Slug     string
	Frontmatter
	Number int
}

type TemplateData struct {
	Sections       map[string]bool
	Title          string
	Author         string
	Date           string
	Summary        string
	PR             string
	TargetServices string
	Motivation     string
	Implementation string
	Metrics        string
	Drawbacks      string
	Alternatives   string
	Impact         string
	Unresolved     string
	Conclusion     string
}

func ParseFile(path string) (*RFC, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(string(data), path)
}

func Parse(content, filePath string) (*RFC, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid RFC format: missing YAML frontmatter delimiters")
	}

	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}

	r := &RFC{
		Frontmatter: fm,
		Body:        strings.TrimSpace(parts[2]),
		FilePath:    filePath,
	}
	return r, nil
}

func (r *RFC) Write() error {
	fmBytes, err := yaml.Marshal(&r.Frontmatter)
	if err != nil {
		return fmt.Errorf("marshal frontmatter: %w", err)
	}

	content := "---\n" + string(fmBytes) + "---\n\n" + r.Body + "\n"
	return os.WriteFile(r.FilePath, []byte(content), 0o644)
}
