package cmd

import (
	"fmt"
	"io/fs"

	"github.com/flaticols/mkrfc/internal/rfc"
	"github.com/flaticols/mkrfc/internal/tpl"
)

func runTemplates(embeddedFS fs.FS) int {
	templates, err := tpl.ListTemplates(embeddedFS, rfc.TemplateDir)
	if err != nil {
		fmt.Printf("error listing templates: %v\n", err)
		return 1
	}

	for _, t := range templates {
		if t.Builtin {
			fmt.Printf("%s (built-in)\n", t.Name)
		} else {
			fmt.Printf("%s\n", t.Name)
		}
	}
	return 0
}
