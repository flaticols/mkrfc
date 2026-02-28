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
  mkrfc new         Create a new RFC
  mkrfc list        List all RFCs
  mkrfc resolve     Resolve (accept/reject) a draft RFC
  mkrfc review      Review a draft RFC (add reviewer, approve, or reject)
  mkrfc templates   List available templates
  mkrfc version     Print version
  mkrfc llm         Print compact LLM-friendly help in XML format
  mkrfc help        Show this help

RFCs are stored in docs/rfc/
Custom templates go in docs/rfc/.templates/
`)
}
