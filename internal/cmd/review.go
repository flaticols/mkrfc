package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh"
	"github.com/flaticols/mkrfc/internal/rfc"
)

func addReviewer(r *rfc.RFC, reviewer string) {
	if reviewer != "" && !slices.Contains(r.Reviewers, reviewer) {
		r.Reviewers = append(r.Reviewers, reviewer)
	}
}

func runReview() int {
	rfcs, err := rfc.Scan(rfc.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning RFCs: %v\n", err)
		return 1
	}

	drafts := rfc.FilterByStatus(rfcs, rfc.StatusDraft)
	if len(drafts) == 0 {
		fmt.Println("No draft RFCs to review.")
		return 0
	}

	options := make([]huh.Option[string], len(drafts))
	for i, r := range drafts {
		options[i] = huh.NewOption(fmt.Sprintf("%04d: %s", r.Number, r.Title), r.FilePath)
	}

	var selectedPath string
	reviewer := gitAuthor()
	var action string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select RFC to review").
				Options(options...).
				Value(&selectedPath),
			huh.NewInput().
				Title("Your name").
				Value(&reviewer),
			huh.NewSelect[string]().
				Title("Action").
				Options(
					huh.NewOption("Add as reviewer", "review"),
					huh.NewOption("Approve (ACCEPTED)", "accept"),
					huh.NewOption("Reject (REJECTED)", "reject"),
				).
				Value(&action),
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

	addReviewer(r, reviewer)

	switch action {
	case "accept", "reject":
		applyResolution(r, action)
	}

	if err = r.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "error writing RFC: %v\n", err)
		return 1
	}

	fmt.Printf("%s updated\n", selectedPath)
	return 0
}

func collectNote(action string) string {
	var note string
	title := "Approval note"
	if action == "reject" {
		title = "Rejection note"
	}
	form := huh.NewForm(huh.NewGroup(
		huh.NewText().
			Title(title).
			Placeholder("Describe the reasoning...").
			Value(&note),
	))
	_ = form.Run()
	return note
}
