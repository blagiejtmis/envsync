package env

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Shell represents a supported shell type for completion generation.
type Shell string

const (
	ShellBash Shell = "bash"
	ShellZsh  Shell = "zsh"
	ShellFish Shell = "fish"
)

// CompletionOptions controls how shell completion snippets are generated.
type CompletionOptions struct {
	// ExportPrefix adds "export " before each assignment (bash/zsh).
	ExportPrefix bool
	// OnlyKeys emits only key names, not values (useful for autocomplete lists).
	OnlyKeys bool
}

// DefaultCompletionOptions returns sensible defaults.
func DefaultCompletionOptions() CompletionOptions {
	return CompletionOptions{
		ExportPrefix: false,
		OnlyKeys:     false,
	}
}

// FprintCompletion writes a shell-specific eval-able snippet to w.
// For bash/zsh it emits KEY=VALUE pairs (one per line, optionally exported).
// For fish it emits `set -x KEY VALUE` statements.
// Entries are sorted by key for deterministic output.
func FprintCompletion(w io.Writer, entries []Entry, shell Shell, opts CompletionOptions) error {
	sorted := make([]Entry, len(entries))
	copy(sorted, entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, e := range sorted {
		if e.Key == "" {
			continue
		}
		var line string
		switch shell {
		case ShellFish:
			if opts.OnlyKeys {
				line = e.Key
			} else {
				line = fmt.Sprintf("set -x %s %s", e.Key, shellQuote(e.Value))
			}
		default: // bash, zsh
			if opts.OnlyKeys {
				line = e.Key
			} else if opts.ExportPrefix {
				line = fmt.Sprintf("export %s=%s", e.Key, shellQuote(e.Value))
			} else {
				line = fmt.Sprintf("%s=%s", e.Key, shellQuote(e.Value))
			}
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

// shellQuote wraps value in single quotes if it contains special characters.
func shellQuote(v string) string {
	if v == "" {
		return "\'\'"
	}
	specials := " \t\n$`\\\"|&;()<>{}!#"
	if strings.ContainsAny(v, specials) {
		escaped := strings.ReplaceAll(v, "'", "'\"'\"'")
		return "'" + escaped + "'"
	}
	return v
}
