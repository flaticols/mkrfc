package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/flaticols/mkrfc/internal/rfc"
)

func handleNoArgs(embeddedFS fs.FS) int {
	rfcs, err := rfc.Scan(rfc.Dir)
	if err != nil || len(rfcs) == 0 {
		return runNew(embeddedFS)
	}
	return runList()
}

func runInfoCmd(cmd string) bool {
	switch cmd {
	case "version":
		printVersion()
	case "llm":
		printLLMHelp()
	case "help", "--help", "-h":
		printUsage()
	default:
		return false
	}
	return true
}

func Run(embeddedFS fs.FS) int {
	if len(os.Args) < 2 {
		return handleNoArgs(embeddedFS)
	}
	if runInfoCmd(os.Args[1]) {
		return 0
	}
	switch os.Args[1] {
	case "new":
		return runNew(embeddedFS)
	case "list":
		return runList()
	case "templates":
		return runTemplates(embeddedFS)
	case "resolve":
		return runResolve()
	case "review", "approve":
		return runReview()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Print(`mkrfc — RFC management tool

Usage:
  mkrfc              List RFCs if any exist, otherwise create a new one
  mkrfc new          Create a new RFC (interactive, or use flags in non-TTY)
  mkrfc list         List all RFCs
  mkrfc resolve      Resolve a draft RFC (accept / reject / supersede)
  mkrfc review       Add reviewer; optionally approve or reject
  mkrfc approve      Alias for review
  mkrfc templates    List available templates
  mkrfc version      Print version
  mkrfc llm          Print compact LLM-friendly help in XML format
  mkrfc help         Show this help

Non-interactive flags (when stdout is not a TTY):
  mkrfc new     --title TEXT [--summary TEXT] [--template NAME] [--sections LIST]
  mkrfc resolve --rfc N --status STATUS [--note TEXT]
  mkrfc review  --rfc N [--action review|accept|reject] [--reviewer NAME] [--note TEXT]

RFCs are stored in docs/rfc/
Custom templates go in docs/rfc/.templates/
`)
}
