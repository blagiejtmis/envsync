package env

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

// FprintHistory writes a human-readable summary of history entries to w.
func FprintHistory(w io.Writer, hf HistoryFile) {
	if len(hf.Entries) == 0 {
		fmt.Fprintln(w, "(no history recorded)")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "#\tTIMESTAMP\tLABEL\tKEYS")
	for i, e := range hf.Entries {
		label := e.Label
		if label == "" {
			label = "-"
		}
		fmt.Fprintf(tw, "%d\t%s\t%s\t%d\n",
			i+1,
			e.Timestamp.Format(time.RFC3339),
			label,
			len(e.Entries),
		)
	}
	tw.Flush()
}

// FprintHistoryDiff writes the diff between two consecutive history entries.
func FprintHistoryDiff(w io.Writer, hf HistoryFile, index int) error {
	if index < 1 || index >= len(hf.Entries) {
		return fmt.Errorf("history: index %d out of range (have %d entries)", index, len(hf.Entries))
	}
	prev := hf.Entries[index-1]
	curr := hf.Entries[index]

	changes := Diff(prev.Entries, curr.Entries)
	if !HasChanges(changes) {
		fmt.Fprintln(w, "(no changes)")
		return nil
	}

	fmt.Fprintf(w, "diff between entry %d and %d\n", index, index+1)
	fmt.Fprintf(w, "  from: %s\n", prev.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "  to:   %s\n\n", curr.Timestamp.Format(time.RFC3339))

	for _, c := range changes {
		switch c.Kind {
		case DiffAdded:
			fmt.Fprintf(w, "+ %s\n", c.Key)
		case DiffRemoved:
			fmt.Fprintf(w, "- %s\n", c.Key)
		case DiffUpdated:
			fmt.Fprintf(w, "~ %s\n", c.Key)
		}
	}
	return nil
}

// HistorySummary returns a one-line summary string for the history file.
func HistorySummary(hf HistoryFile) string {
	if len(hf.Entries) == 0 {
		return "no history"
	}
	last := hf.Entries[len(hf.Entries)-1]
	parts := []string{
		fmt.Sprintf("%d snapshot(s)", len(hf.Entries)),
		fmt.Sprintf("latest: %s", last.Timestamp.Format(time.RFC3339)),
	}
	if last.Label != "" {
		parts = append(parts, fmt.Sprintf("label: %s", last.Label))
	}
	return strings.Join(parts, " | ")
}
