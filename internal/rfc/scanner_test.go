package rfc_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/flaticols/mkrfc/internal/rfc"
)

func writeFixture(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const fixture1 = `---
title: "First RFC"
authors:
  - "Alice"
status: "DRAFT"
date: "2024-01-01"
tags: []
---

Body text.
`

const fixture2 = `---
title: "Second RFC"
authors:
  - "Bob"
status: "ACCEPTED"
date: "2024-02-01"
tags: []
---

Body text.
`

func TestScan(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "0001-first-rfc.md", fixture1)
	writeFixture(t, dir, "0002-second-rfc.md", fixture2)

	rfcs, err := rfc.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rfcs) != 2 {
		t.Fatalf("expected 2 RFCs, got %d", len(rfcs))
	}
	if rfcs[0].Number != 1 {
		t.Errorf("expected first RFC number 1, got %d", rfcs[0].Number)
	}
	if rfcs[1].Number != 2 {
		t.Errorf("expected second RFC number 2, got %d", rfcs[1].Number)
	}
	if rfcs[0].Title != "First RFC" {
		t.Errorf("unexpected title: %s", rfcs[0].Title)
	}
}

func TestNextNumber(t *testing.T) {
	dir := t.TempDir()

	n, err := rfc.NextNumber(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("expected 1 for empty dir, got %d", n)
	}

	writeFixture(t, dir, "0001-first.md", fixture1)
	writeFixture(t, dir, "0003-third.md", fixture2)

	n, err = rfc.NextNumber(dir)
	if err != nil {
		t.Fatal(err)
	}
	if n != 4 {
		t.Errorf("expected 4, got %d", n)
	}
}

func TestFilterByStatus(t *testing.T) {
	dir := t.TempDir()
	writeFixture(t, dir, "0001-first.md", fixture1)
	writeFixture(t, dir, "0002-second.md", fixture2)

	rfcs, err := rfc.Scan(dir)
	if err != nil {
		t.Fatal(err)
	}

	drafts := rfc.FilterByStatus(rfcs, rfc.StatusDraft)
	if len(drafts) != 1 {
		t.Errorf("expected 1 draft, got %d", len(drafts))
	}
	accepted := rfc.FilterByStatus(rfcs, rfc.StatusAccepted)
	if len(accepted) != 1 {
		t.Errorf("expected 1 accepted, got %d", len(accepted))
	}
}
