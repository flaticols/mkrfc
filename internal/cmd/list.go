package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/flaticols/mkrfc/internal/rfc"
)

func runList() int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	statusFlag := fs.String("status", "", "filter by status (DRAFT, ACCEPTED, REJECTED, SUPERSEDED)")
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}
	if err := fs.Parse(args); err != nil {
		return 1
	}

	rfcs, err := rfc.Scan(rfc.Dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scanning RFCs: %v\n", err)
		return 1
	}

	if *statusFlag != "" {
		rfcs = rfc.FilterByStatus(rfcs, rfc.Status(strings.ToUpper(*statusFlag)))
	}

	if len(rfcs) == 0 {
		fmt.Println("No RFCs found.")
		return 0
	}

	fmt.Printf("%-6s  %-40s  %-12s  %-12s  %s\n", "NUM", "TITLE", "STATUS", "DATE", "AUTHORS")
	fmt.Println(strings.Repeat("-", 90))
	for _, r := range rfcs {
		author := ""
		if len(r.Authors) > 0 {
			author = r.Authors[0]
		}
		title := r.Title
		if len(title) > 40 {
			title = title[:37] + "..."
		}
		fmt.Printf("%04d    %-40s  %-12s  %-12s  %s\n", r.Number, title, r.Status, r.Date, author)
	}
	return 0
}
