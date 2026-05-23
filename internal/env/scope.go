package env

import "fmt"

// Scope represents a named environment tier (e.g. development, staging, production).
type Scope struct {
	Name     string
	Priority int // higher priority wins during merge
}

// ScopeOptions controls how ScopeFilter behaves.
type ScopeOptions struct {
	// AllowedScopes is the ordered list of known scope names (lowest to highest priority).
	AllowedScopes []string
	// AnnotationKey is the comment annotation used to tag an entry with a scope, e.g. "scope".
	AnnotationKey string
}

// DefaultScopeOptions returns sensible defaults.
func DefaultScopeOptions() ScopeOptions {
	return ScopeOptions{
		AllowedScopes: []string{"development", "staging", "production"},
		AnnotationKey: "scope",
	}
}

// ScopeFilter returns only entries whose scope annotation matches one of the
// provided target scopes. Entries with no scope annotation are always included.
func ScopeFilter(entries []Entry, targets []string, opts ScopeOptions) ([]Entry, error) {
	if len(targets) == 0 {
		return Clone(entries), nil
	}

	allowed := make(map[string]bool, len(opts.AllowedScopes))
	for _, s := range opts.AllowedScopes {
		allowed[s] = true
	}
	for _, t := range targets {
		if !allowed[t] {
			return nil, fmt.Errorf("env/scope: unknown scope %q", t)
		}
	}

	targetSet := make(map[string]bool, len(targets))
	for _, t := range targets {
		targetSet[t] = true
	}

	var out []Entry
	for _, e := range entries {
		scope := extractAnnotation(e.Comment, opts.AnnotationKey)
		if scope == "" || targetSet[scope] {
			out = append(out, e)
		}
	}
	return out, nil
}

// extractAnnotation parses a comment string for a key=value annotation.
// e.g. "# scope=production" with key "scope" returns "production".
func extractAnnotation(comment, key string) string {
	if comment == "" || key == "" {
		return ""
	}
	prefix := key + "="
	for _, part := range splitWords(comment) {
		if len(part) > len(prefix) && part[:len(prefix)] == prefix {
			return part[len(prefix):]
		}
	}
	return ""
}

// splitWords splits a comment string (stripping leading '#') into words.
func splitWords(s string) []string {
	// strip leading '#' and spaces
	for len(s) > 0 && (s[0] == '#' || s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	var words []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				words = append(words, s[start:i])
				start = -1
			}
		} else {
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		words = append(words, s[start:])
	}
	return words
}
