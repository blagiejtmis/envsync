package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// PromoteResult holds the outcome of a promotion operation for reporting.
type PromoteResult struct {
	FromEnv  string
	ToEnv    string
	Promoted []string
	Skipped  []string
}

// FprintPromoteResult writes a human-readable summary of a promotion to w.
func FprintPromoteResult(w io.Writer, r PromoteResult) {
	from := r.FromEnv
	if from == "" {
		from = "source"
	}
	to := r.ToEnv
	if to == "" {
		to = "destination"
	}

	fmt.Fprintf(w, "Promote: %s → %s\n", from, to)

	if len(r.Promoted) == 0 {
		fmt.Fprintln(w, "  No keys promoted.")
	} else {
		sorted := append([]string(nil), r.Promoted...)
		sort.Strings(sorted)
		fmt.Fprintf(w, "  Promoted (%d):\n", len(sorted))
		for _, k := range sorted {
			fmt.Fprintf(w, "    + %s\n", k)
		}
	}

	if len(r.Skipped) > 0 {
		sorted := append([]string(nil), r.Skipped...)
		sort.Strings(sorted)
		fmt.Fprintf(w, "  Skipped (%d, already exist):\n", len(sorted))
		for _, k := range sorted {
			fmt.Fprintf(w, "    ~ %s\n", k)
		}
	}
}

// PromoteSummary returns a single-line summary string.
func PromoteSummary(r PromoteResult) string {
	parts := []string{}
	if len(r.Promoted) > 0 {
		parts = append(parts, fmt.Sprintf("%d promoted", len(r.Promoted)))
	}
	if len(r.Skipped) > 0 {
		parts = append(parts, fmt.Sprintf("%d skipped", len(r.Skipped)))
	}
	if len(parts) == 0 {
		return "nothing to promote"
	}
	return strings.Join(parts, ", ")
}
