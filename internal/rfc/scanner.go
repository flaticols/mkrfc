package rfc

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	Dir         = "docs/rfc"
	TemplateDir = "docs/rfc/.templates"
)

func Scan(dir string) ([]*RFC, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "[0-9][0-9][0-9][0-9]-*.md"))
	if err != nil {
		return nil, err
	}

	var rfcs []*RFC
	for _, path := range entries {
		r, err := ParseFile(path)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		r.Number, r.Slug = parseFilename(filepath.Base(path))
		rfcs = append(rfcs, r)
	}

	sort.Slice(rfcs, func(i, j int) bool {
		return rfcs[i].Number < rfcs[j].Number
	})
	return rfcs, nil
}

func NextNumber(dir string) (int, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 1, nil
	}
	rfcs, err := Scan(dir)
	if err != nil {
		return 0, err
	}
	if len(rfcs) == 0 {
		return 1, nil
	}
	return rfcs[len(rfcs)-1].Number + 1, nil
}

func FilterByStatus(rfcs []*RFC, status Status) []*RFC {
	var out []*RFC
	for _, r := range rfcs {
		if r.Status == status {
			out = append(out, r)
		}
	}
	return out
}

func parseFilename(name string) (int, string) {
	// name is like: 0001-my-rfc-slug.md
	base := strings.TrimSuffix(name, ".md")
	parts := strings.SplitN(base, "-", 2)
	if len(parts) != 2 {
		return 0, base
	}
	n, _ := strconv.Atoi(parts[0])
	return n, parts[1]
}
