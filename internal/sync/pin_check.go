// Package sync provides syncing helpers for envsync.
package sync

import (
	"fmt"
	"io"
	"sort"

	"github.com/example/envsync/internal/env"
)

// PinViolation describes a key whose current value differs from its pinned value.
type PinViolation struct {
	Key      string
	Pinned   string
	Current  string
}

// CheckPins compares entries against pinned values and returns any violations.
// A violation occurs when a pinned key exists in entries but its value has changed.
func CheckPins(entries []env.Entry, pf env.PinFile) []PinViolation {
	pinned := env.PinnedKeys(pf)
	index := make(map[string]string, len(entries))
	for _, e := range entries {
		index[e.Key] = e.Value
	}

	var violations []PinViolation
	for key, pin := range pinned {
		current, ok := index[key]
		if !ok {
			// Key removed entirely — still a violation.
			violations = append(violations, PinViolation{
				Key:     key,
				Pinned:  pin.Value,
				Current: "<missing>",
			})
			continue
		}
		if current != pin.Value {
			violations = append(violations, PinViolation{
				Key:     key,
				Pinned:  pin.Value,
				Current: current,
			})
		}
	}
	sort.Slice(violations, func(i, j int) bool {
		return violations[i].Key < violations[j].Key
	})
	return violations
}

// FprintPinResult writes a human-readable pin check result to w.
func FprintPinResult(w io.Writer, violations []PinViolation) {
	if len(violations) == 0 {
		fmt.Fprintln(w, "pin check: all pinned values match")
		return
	}
	fmt.Fprintf(w, "pin check: %d violation(s):\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(w, "  %-24s pinned=%q current=%q\n", v.Key, v.Pinned, v.Current)
	}
}
