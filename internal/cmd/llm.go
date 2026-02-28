package cmd

import "fmt"

func printLLMHelp() {
	fmt.Print(`<mkrfc>
<desc>CLI tool for creating and managing project RFCs. Stores RFCs as Markdown files with YAML frontmatter in docs/rfc/. Interactive TUI forms powered by charmbracelet/huh.</desc>
<commands>
mkrfc                  if no RFCs exist, run new; otherwise run list
mkrfc new              create RFC interactively (template, title, summary, sections)
mkrfc list             list all RFCs: NUM, TITLE, STATUS, DATE, AUTHORS
mkrfc list --status S  filter by status (DRAFT, ACCEPTED, REJECTED, SUPERSEDED)
mkrfc resolve          transition a draft RFC: pick RFC, new status, resolution note
mkrfc review           add reviewer to a draft; optionally approve or reject
mkrfc approve          alias for review
mkrfc templates        list available templates (built-in and custom)
mkrfc version          print installed version
mkrfc llm              print this help
mkrfc help             show usage
</commands>
<storage>
RFCs stored in: docs/rfc/NNNN-slug.md  (e.g. 0001-add-async-job-queue.md)
Custom templates: docs/rfc/.templates/*.md
Default template: embedded in binary via //go:embed
</storage>
<frontmatter>
title:         string
authors:       []string
reviewers:     []string
status:        DRAFT | ACCEPTED | REJECTED | SUPERSEDED
date:          YYYY-MM-DD
resolved_date: YYYY-MM-DD  (set by resolve/review)
pr:            URL          (auto-detected from gh CLI, omitted if empty)
tags:          []string
</frontmatter>
<template-data>
Title          string  RFC title
Author         string  git config user.name + user.email
Date           string  creation date YYYY-MM-DD
Summary        string  one-line summary from form
PR             string  PR URL from gh CLI (empty if none)
Sections       map[string]bool  selected optional sections; check with: {{ if index .Sections "key" }}
TargetServices string  optional, empty by default
Motivation     string  optional, empty by default
Implementation string  optional, empty by default
Metrics        string  optional, empty by default
Drawbacks      string  optional, empty by default
Alternatives   string  optional, empty by default
Impact         string  optional, empty by default
Unresolved     string  optional, empty by default
Conclusion     string  optional, empty by default
</template-data>
<optional-sections>
Sections pre-selected by default: motivation, implementation, drawbacks, alternatives
All available keys: motivation, implementation, metrics, drawbacks, alternatives, impact, unresolved, conclusion
Unselected sections are omitted entirely from the output.
</optional-sections>
<pr-detection>
If gh CLI is installed and current branch has an open PR:
  pr field in frontmatter is set to the PR URL
  RFC header table shows a "View PR" link
Command: gh pr view --json url -q .url
</pr-detection>
<review-flow>
1. Pick draft RFC
2. Enter reviewer name (pre-filled from git config)
3. Action: "Add as reviewer" | "Approve (ACCEPTED)" | "Reject (REJECTED)"
4. If approve/reject: enter note → appended as ## Resolution section
</review-flow>
<non-tty>
When stdin is not a TTY (piped, CI), interactive forms are skipped and flags are used instead.

mkrfc new --title TEXT [--summary TEXT] [--template NAME] [--sections LIST]
  --title      RFC title (required)
  --summary    brief summary (default: "")
  --template   template name (default: first available)
  --sections   comma-separated section keys (default: motivation,implementation,drawbacks,alternatives)

mkrfc resolve --rfc N --status STATUS [--note TEXT]
  --rfc        RFC number (required)
  --status     ACCEPTED | REJECTED | SUPERSEDED (required)
  --note       resolution note

mkrfc review --rfc N [--reviewer NAME] [--action ACTION] [--note TEXT]
  --rfc        RFC number (required)
  --reviewer   reviewer name (default: git config user)
  --action     review | accept | reject (default: review)
  --note       resolution note (for accept/reject)
</non-tty>
<examples>
mkrfc                              # new RFC if none exist
mkrfc new                          # create RFC interactively
mkrfc list                         # list all RFCs
mkrfc list --status DRAFT          # show only drafts
mkrfc resolve                      # resolve a draft RFC
mkrfc review                       # add reviewer or approve/reject
mkrfc templates                    # list available templates
</examples>
</mkrfc>
`)
}
