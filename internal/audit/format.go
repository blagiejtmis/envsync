package audit

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

// Fprint writes a human-readable table of audit entries to w.
func Fprint(w io.Writer, entries []Entry) error {
	if len(entries) == 0 {
		_, err := fmt.Fprintln(w, "(no audit entries)")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tKIND\tKEY\tUSER\tMESSAGE")
	fmt.Fprintln(tw, strings.Repeat("-", 9)+"\t"+strings.Repeat("-", 6)+"\t"+strings.Repeat("-", 3)+"\t"+strings.Repeat("-", 4)+"\t"+strings.Repeat("-", 7))

	for _, e := range entries {
		ts := e.Timestamp.In(time.Local).Format(timeFormat)
		user := e.User
		if user == "" {
			user = "-"
		}
		msg := e.Message
		if msg == "" {
			msg = "-"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", ts, e.Kind, e.Key, user, msg)
	}
	return tw.Flush()
}

// Summary returns a one-line string summarising an Entry.
func Summary(e Entry) string {
	ts := e.Timestamp.UTC().Format(timeFormat)
	return fmt.Sprintf("[%s] %s key=%q user=%q", ts, e.Kind, e.Key, e.User)
}
