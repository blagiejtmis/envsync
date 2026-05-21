package version

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Fprint writes a formatted version table to w.
// Each field is printed on its own line for easy parsing by scripts.
func Fprint(w io.Writer, i Info) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Version:\t%s\n", i.String())
	fmt.Fprintf(tw, "Commit:\t%s\n", i.Commit)
	fmt.Fprintf(tw, "Built:\t%s\n", i.BuildDate)
	tw.Flush()
}

// FprintFull writes a single-line full version string to w followed by a
// newline. Suitable for --version flag output.
func FprintFull(w io.Writer, i Info) {
	fmt.Fprintln(w, i.Full())
}
