package env

import (
	"fmt"
	"io"
)

// FprintCompare writes a human-readable comparison report to w.
func FprintCompare(w io.Writer, r CompareResult) {
	if !r.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	for _, e := range r.OnlyInA {
		fmt.Fprintf(w, "- %s\n", e.Key)
	}
	for _, e := range r.OnlyInB {
		fmt.Fprintf(w, "+ %s\n", e.Key)
	}
	for _, p := range r.Changed {
		fmt.Fprintf(w, "~ %s\n", p.A.Key)
	}
}

// FprintCompareVerbose writes a detailed comparison report including values to w.
// Sensitive keys are redacted.
func FprintCompareVerbose(w io.Writer, r CompareResult) {
	if !r.HasDifferences() {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	if len(r.OnlyInA) > 0 {
		fmt.Fprintln(w, "Only in A:")
		for _, e := range r.OnlyInA {
			v := e.Value
			if IsSensitive(e.Key) {
				v = "[redacted]"
			}
			fmt.Fprintf(w, "  - %s=%s\n", e.Key, v)
		}
	}

	if len(r.OnlyInB) > 0 {
		fmt.Fprintln(w, "Only in B:")
		for _, e := range r.OnlyInB {
			v := e.Value
			if IsSensitive(e.Key) {
				v = "[redacted]"
			}
			fmt.Fprintf(w, "  + %s=%s\n", e.Key, v)
		}
	}

	if len(r.Changed) > 0 {
		fmt.Fprintln(w, "Changed:")
		for _, p := range r.Changed {
			if IsSensitive(p.A.Key) {
				fmt.Fprintf(w, "  ~ %s: [redacted] -> [redacted]\n", p.A.Key)
			} else {
				fmt.Fprintf(w, "  ~ %s: %q -> %q\n", p.A.Key, p.A.Value, p.B.Value)
			}
		}
	}
}

// CompareSummary returns a one-line summary of the comparison result.
func CompareSummary(r CompareResult) string {
	if !r.HasDifferences() {
		return fmt.Sprintf("identical (%d keys)", len(r.Identical))
	}
	return fmt.Sprintf("%d added, %d removed, %d changed, %d identical",
		len(r.OnlyInB), len(r.OnlyInA), len(r.Changed), len(r.Identical))
}
