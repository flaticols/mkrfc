package main

import (
	"embed"
	"os"

	"github.com/flaticols/mkrfc/internal/cmd"
)

//go:embed templates/default.md
var embeddedTemplates embed.FS

func main() {
	os.Exit(cmd.Run(embeddedTemplates))
}
