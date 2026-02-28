package cmd

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/flaticols/mkrfc/internal/rfc"
)

type reviewInput struct {
	path, reviewer, action, note string
}

func buildReviewForm(drafts []*rfc.RFC) (reviewInput, error) {
	options := make([]huh.Option[string], len(drafts))
	for i, r := range drafts {
		options[i] = huh.NewOption(fmt.Sprintf("%04d: %s", r.Number, r.Title), r.FilePath)
	}

	var in reviewInput
	in.reviewer = gitAuthor()

	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Select RFC to review").
			Options(options...).
			Value(&in.path),
		huh.NewInput().
			Title("Your name").
			Value(&in.reviewer),
		huh.NewSelect[string]().
			Title("Action").
			Options(
				huh.NewOption("Add as reviewer", "review"),
				huh.NewOption("Approve (ACCEPTED)", "accept"),
				huh.NewOption("Reject (REJECTED)", "reject"),
			).
			Value(&in.action),
	))
	if err := form.Run(); err != nil {
		return in, err
	}
	return in, nil
}

func parseReviewFlags(drafts []*rfc.RFC) (reviewInput, error) {
	fs := flag.NewFlagSet("review", flag.ContinueOnError)
	rfcNum := fs.Int("rfc", 0, "RFC number (required)")
	reviewer := fs.String("reviewer", gitAuthor(), "reviewer name")
	action := fs.String("action", "review", "action: review, accept, reject")
	note := fs.String("note", "", "resolution note (for accept/reject)")

	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if err := fs.Parse(args); err != nil {
		return reviewInput{}, err
	}
	if *rfcNum == 0 {
		return reviewInput{}, fmt.Errorf("--rfc is required in non-interactive mode")
	}
	for _, r := range drafts {
		if r.Number == *rfcNum {
			return reviewInput{path: r.FilePath, reviewer: *reviewer, action: *action, note: *note}, nil
		}
	}
	return reviewInput{}, fmt.Errorf("RFC %04d not found in DRAFT status", *rfcNum)
}

func getReviewInput(drafts []*rfc.RFC) (reviewInput, error) {
	if isTTY() {
		return buildReviewForm(drafts)
	}
	return parseReviewFlags(drafts)
}

func addReviewer(r *rfc.RFC, reviewer string) {
	if reviewer != "" && !slices.Contains(r.Reviewers, reviewer) {
		r.Reviewers = append(r.Reviewers, reviewer)
	}
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

	in, err := getReviewInput(drafts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	r, err := rfc.ParseFile(in.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing RFC: %v\n", err)
		return 1
	}

	addReviewer(r, in.reviewer)

	if in.action != "review" {
		// For non-interactive, note comes from the flag; for interactive,
		// applyResolution collects it via a second form (TTY only).
		if isTTY() {
			applyResolution(r, in.action)
		} else {
			applyResolutionWithNote(r, in.action, in.note)
		}
	}

	if err = r.Write(); err != nil {
		fmt.Fprintf(os.Stderr, "error writing RFC: %v\n", err)
		return 1
	}

	fmt.Printf("%s updated\n", in.path)
	return 0
}

func applyResolutionWithNote(r *rfc.RFC, action, note string) {
	status := rfc.StatusAccepted
	if action == "reject" {
		status = rfc.StatusRejected
	}
	today := time.Now().Format("2006-01-02")
	r.Status = status
	r.ResolvedDate = today
	resolution := fmt.Sprintf("\n\n---\n\n## Resolution\n\n**Status:** %s  \n**Date:** %s\n\n%s",
		string(status), today, strings.TrimSpace(note))
	r.Body += resolution
}
