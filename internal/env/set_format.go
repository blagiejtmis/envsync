package env

import (
	"fmt"
	"io"
	"sort"
)

// SetResult records what happened during a Set operation.
type SetResult struct {
	Added    []string // keys that were newly added
	Updated  []string // keys whose values changed
	Unchanged []string // keys whose values were already equal
}

// SetDiff applies pairs to entries and returns both the updated entries and a
// SetResult describing what changed. The original slice is not mutated.
func SetDiff(entries []Entry, pairs map[string]string, opts SetOptions) ([]Entry, SetResult, error) {
	var result SetResult

	for key, value := range pairs {
		idx := indexOfEntry(entries, key)
		if idx < 0 {
			result.Added = append(result.Added, key)
		} else if entries[idx].Value == value {
			result.Unchanged = append(result.Unchanged, key)
		} else {
			result.Updated = append(result.Updated, key)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Updated)
	sort.Strings(result.Unchanged)

	out, err := Set(entries, pairs, opts)
	if err != nil {
		return nil, SetResult{}, err
	}
	return out, result, nil
}

// FprintSetResult writes a human-readable summary of a SetResult to w.
func FprintSetResult(w io.Writer, r SetResult) {
	if len(r.Added) == 0 && len(r.Updated) == 0 && len(r.Unchanged) == 0 {
		fmt.Fprintln(w, "no changes")
		return
	}
	for _, k := range r.Added {
		fmt.Fprintf(w, "+ %s (added)\n", k)
	}
	for _, k := range r.Updated {
		fmt.Fprintf(w, "~ %s (updated)\n", k)
	}
	for _, k := range r.Unchanged {
		fmt.Fprintf(w, "  %s (unchanged)\n", k)
	}
}

// SetSummary returns a one-line summary string for a SetResult.
func SetSummary(r SetResult) string {
	return fmt.Sprintf("%d added, %d updated, %d unchanged",
		len(r.Added), len(r.Updated), len(r.Unchanged))
}
