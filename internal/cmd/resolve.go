package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/flaticols/mkrfc/internal/rfc"
)

type resolveInput struct {
	path, status, note string
}

func buildResolveForm(drafts []*rfc.RFC) (resolveInput, error) {
	options := make([]huh.Option[string], len(drafts))
	for i, r := range drafts {
		options[i] = huh.NewOption(fmt.Sprintf("%04d: %s", r.Number, r.Title), r.FilePath)
	}

	var in resolveInput
	statusOptions := []huh.Option[string]{
		huh.NewOption("ACCEPTED", string(rfc.StatusAccepted)),
		huh.NewOption("REJECTED", string(rfc.StatusRejected)),
		huh.NewOption("SUPERSEDED", string(rfc.StatusSuperseded)),
	}

	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().Title("Select RFC to resolve").Options(options...).Value(&in.path),
		huh.NewSelect[string]().Title("New status").Options(statusOptions...).Value(&in.status),
		huh.NewText().Title("Resolution note").Placeholder("Describe why this RFC was accepted/rejected...").Value(&in.note),
	))
	if err := form.Run(); err != nil {
		return in, err
	}
	return in, nil
}

func parseResolveFlags(drafts []*rfc.RFC) (resolveInput, error) {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	rfcNum := fs.Int("rfc", 0, "RFC number (required)")
	status := fs.String("status", "", "new status: ACCEPTED, REJECTED, SUPERSEDED (required)")
	note := fs.String("note", "", "resolution note")

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if err := fs.Parse(args); err != nil {
		return resolveInput{}, err
	}
	if *rfcNum == 0 {
		return resolveInput{}, fmt.Errorf("--rfc is required in non-interactive mode")
	}
	if *status == "" {
		return resolveInput{}, fmt.Errorf("--status is required in non-interactive mode")
	}
	for _, r := range drafts {
		if r.Number == *rfcNum {
			return resolveInput{path: r.FilePath, status: strings.ToUpper(*status), note: *note}, nil
		}
	}
	return resolveInput{}, fmt.Errorf("RFC %04d not found in DRAFT status", *rfcNum)
}

func getResolveInput(drafts []*rfc.RFC) (resolveInput, error) {
	if isTTY() {
		return buildResolveForm(drafts)
	}
	return parseResolveFlags(drafts)
}

func applyResolution(r *rfc.RFC, action string) {
	status := rfc.StatusAccepted
	if action == "reject" {
		status = rfc.StatusRejected
	}
	note := collectNote(action)
	today := time.Now().Format("2006-01-02")
	r.Status = status
	r.ResolvedDate = today
	resolution := fmt.Sprintf("\n\n---\n\n## Resolution\n\n**Status:** %s  \n**Date:** %s\n\n%s",
		string(status), today, strings.TrimSpace(note))
	r.Body += resolution
}

func runResolve() int {
	rfcs, err := rfc.Scan(rfc.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning RFCs: %v\n", err)
		return 1
	}

	drafts := rfc.FilterByStatus(rfcs, rfc.StatusDraft)
	if len(drafts) == 0 {
		fmt.Println("No draft RFCs to resolve.")
		return 0
	}

	in, err := getResolveInput(drafts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	r, err := rfc.ParseFile(in.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing RFC: %v\n", err)
		return 1
	}

	today := time.Now().Format("2006-01-02")
	r.Status = rfc.Status(in.status)
	r.ResolvedDate = today
	resolution := fmt.Sprintf("\n\n---\n\n## Resolution\n\n**Status:** %s  \n**Date:** %s\n\n%s",
		in.status, today, strings.TrimSpace(in.note))
	r.Body += resolution

	if err = r.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "error writing RFC: %v\n", err)
		return 1
	}

	fmt.Printf("%s resolved as %s\n", in.path, in.status)
	return 0
}
