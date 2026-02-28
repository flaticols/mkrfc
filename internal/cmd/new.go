package cmd

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/flaticols/mkrfc/internal/rfc"
	"github.com/flaticols/mkrfc/internal/tpl"
)

type newFormResult struct {
	title            string
	summary          string
	selectedTemplate string
	sections         []string
}

func sectionOptions() []huh.Option[string] {
	return []huh.Option[string]{
		huh.NewOption("Motivation", "motivation").Selected(true),
		huh.NewOption("Proposed Implementation", "implementation").Selected(true),
		huh.NewOption("Metrics & Dashboards", "metrics"),
		huh.NewOption("Drawbacks", "drawbacks").Selected(true),
		huh.NewOption("Alternatives", "alternatives").Selected(true),
		huh.NewOption("Potential Impact & Dependencies", "impact"),
		huh.NewOption("Unresolved Questions", "unresolved"),
		huh.NewOption("Conclusion", "conclusion"),
	}
}

func buildNewForm(templates []tpl.Template) (newFormResult, error) {
	var result newFormResult

	var fields []huh.Field

	if len(templates) > 1 {
		options := make([]huh.Option[string], len(templates))
		for i, t := range templates {
			label := t.Name
			if t.Builtin {
				label += " (built-in)"
			}
			options[i] = huh.NewOption(label, t.Name)
		}
		fields = append(fields, huh.NewSelect[string]().
			Title("Template").
			Options(options...).
			Value(&result.selectedTemplate))
	} else {
		result.selectedTemplate = templates[0].Name
	}

	fields = append(fields,
		huh.NewInput().
			Title("RFC Title").
			Placeholder("My Proposal").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("title is required")
				}
				return nil
			}).
			Value(&result.title),
		huh.NewText().
			Title("Summary").
			Placeholder("Brief description of the proposal...").
			Value(&result.summary),
		huh.NewMultiSelect[string]().
			Title("Optional Sections").
			Description("Choose which sections to include").
			Options(sectionOptions()...).
			Value(&result.sections),
	)

	form := huh.NewForm(huh.NewGroup(fields...))
	if err := form.Run(); err != nil {
		return result, err
	}
	return result, nil
}

func parseNewFlags(templates []tpl.Template) (newFormResult, error) {
	fs := flag.NewFlagSet("new", flag.ContinueOnError)
	title := fs.String("title", "", "RFC title (required)")
	summary := fs.String("summary", "", "brief summary")
	tmplName := fs.String("template", templates[0].Name, "template name")
	sectionsRaw := fs.String("sections", "motivation,implementation,drawbacks,alternatives", "comma-separated sections to include")

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if err := fs.Parse(args); err != nil {
		return newFormResult{}, err
	}
	if strings.TrimSpace(*title) == "" {
		return newFormResult{}, fmt.Errorf("--title is required in non-interactive mode")
	}

	var sections []string
	for _, s := range strings.Split(*sectionsRaw, ",") {
		if s = strings.TrimSpace(s); s != "" {
			sections = append(sections, s)
		}
	}

	return newFormResult{
		title:            *title,
		summary:          *summary,
		selectedTemplate: *tmplName,
		sections:         sections,
	}, nil
}

func getNewInput(templates []tpl.Template) (newFormResult, error) {
	if isTTY() {
		return buildNewForm(templates)
	}
	return parseNewFlags(templates)
}

func templateContent(templates []tpl.Template, name string) string {
	for _, t := range templates {
		if t.Name == name {
			return t.Content
		}
	}
	return templates[0].Content
}

func runNew(embeddedFS fs.FS) int {
	templates, err := tpl.ListTemplates(embeddedFS, rfc.TemplateDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error listing templates: %v\n", err)
		return 1
	}

	result, err := getNewInput(templates)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	author := gitAuthor()
	date := time.Now().Format("2006-01-02")

	num, err := rfc.NextNumber(rfc.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error getting next number: %v\n", err)
		return 1
	}

	sections := make(map[string]bool, len(result.sections))
	for _, s := range result.sections {
		sections[s] = true
	}

	slug := slugify(result.title)
	content := templateContent(templates, result.selectedTemplate)

	data := rfc.TemplateData{
		Title:    result.title,
		Author:   author,
		Date:     date,
		Summary:  result.summary,
		Sections: sections,
		PR:       detectPRLink(),
	}

	rendered, err := tpl.Render(content, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error rendering template: %v\n", err)
		return 1
	}

	if err = os.MkdirAll(rfc.Dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error creating directory: %v\n", err)
		return 1
	}

	filename := fmt.Sprintf("%04d-%s.md", num, slug)
	filePath := rfc.Dir + "/" + filename

	if err = os.WriteFile(filePath, []byte(rendered), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		return 1
	}

	fmt.Println(filePath)
	return 0
}

func gitAuthor() string {
	name := gitConfig("user.name")
	email := gitConfig("user.email")
	if name == "" {
		return "Unknown"
	}
	if email == "" {
		return name
	}
	return name + " <" + email + ">"
}

func gitConfig(key string) string {
	out, err := exec.Command("git", "config", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func detectPRLink() string {
	if _, err := exec.LookPath("gh"); err != nil {
		return ""
	}
	out, err := exec.Command("gh", "pr", "view", "--json", "url", "-q", ".url").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(s)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
