package env

import (
	"fmt"
	"strings"
)

// SchemaEntry describes a single expected environment variable.
type SchemaEntry struct {
	Key      string
	Required bool
	Default  string
	Comment  string
}

// Schema is an ordered list of schema entries.
type Schema []SchemaEntry

// ParseSchema reads a .env.schema file (same format as .env but with
// optional inline annotations in comments above each key).
//
// Supported annotations (place on the line directly above the key):
//
//	# @required
//	# @default some_value
func ParseSchema(lines []string) (Schema, error) {
	var schema Schema
	var pendingComment string
	var required bool
	var defaultVal string

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			pendingComment = ""
			required = false
			defaultVal = ""
			continue
		}
		if strings.HasPrefix(line, "#") {
			anno := strings.TrimSpace(strings.TrimPrefix(line, "#"))
			if anno == "@required" {
				required = true
			} else if strings.HasPrefix(anno, "@default ") {
				defaultVal = strings.TrimPrefix(anno, "@default ")
			} else {
				pendingComment = anno
			}
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid schema line: %q", line)
		}
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("empty key in schema line: %q", line)
		}
		schema = append(schema, SchemaEntry{
			Key:      key,
			Required: required,
			Default:  defaultVal,
			Comment:  pendingComment,
		})
		pendingComment = ""
		required = false
		defaultVal = ""
	}
	return schema, nil
}

// ValidateAgainstSchema checks entries against a schema and returns a list
// of human-readable violation messages.
func ValidateAgainstSchema(entries []Entry, schema Schema) []string {
	present := make(map[string]string, len(entries))
	for _, e := range entries {
		present[e.Key] = e.Value
	}
	var violations []string
	for _, s := range schema {
		val, ok := present[s.Key]
		if !ok {
			if s.Required {
				violations = append(violations, fmt.Sprintf("missing required key: %s", s.Key))
			}
			continue
		}
		if s.Required && val == "" {
			violations = append(violations, fmt.Sprintf("required key has blank value: %s", s.Key))
		}
	}
	return violations
}
