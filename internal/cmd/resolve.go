package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/flaticols/mkrfc/internal/rfc"
)

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

	options := make([]huh.Option[string], len(drafts))
	for i, r := range drafts {
		label := fmt.Sprintf("%04d: %s", r.Number, r.Title)
		options[i] = huh.NewOption(label, r.FilePath)
	}

	var selectedPath string
	var newStatus string
	var note string

	statusOptions := []huh.Option[string]{
		huh.NewOption("ACCEPTED", string(rfc.StatusAccepted)),
		huh.NewOption("REJECTED", string(rfc.StatusRejected)),
		huh.NewOption("SUPERSEDED", string(rfc.StatusSuperseded)),
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select RFC to resolve").
				Options(options...).
				Value(&selectedPath),
			huh.NewSelect[string]().
				Title("New status").
				Options(statusOptions...).
				Value(&newStatus),
			huh.NewText().
				Title("Resolution note").
				Placeholder("Describe why this RFC was accepted/rejected...").
				Value(&note),
		),
	)

	if err = form.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "form error: %v\n", err)
		return 1
	}

	r, err := rfc.ParseFile(selectedPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing RFC: %v\n", err)
		return 1
	}

	today := time.Now().Format("2006-01-02")
	r.Status = rfc.Status(newStatus)
	r.ResolvedDate = today

	resolution := fmt.Sprintf("\n\n---\n\n## Resolution\n\n**Status:** %s  \n**Date:** %s\n\n%s",
		newStatus, today, strings.TrimSpace(note))
	r.Body += resolution

	if err = r.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "error writing RFC: %v\n", err)
		return 1
	}

	fmt.Printf("RFC %s resolved as %s\n", selectedPath, newStatus)
	return 0
}
