package env

import (
	"testing"
)

func TestIsSensitive(t *testing.T) {
	cases := []struct {
		key       string
		wantSensitive bool
	}{
		{"DB_PASSWORD", true},
		{"API_KEY", true},
		{"GITHUB_TOKEN", true},
		{"AWS_SECRET_ACCESS_KEY", true},
		{"PRIVATE_KEY_PATH", true},
		{"AUTH_HEADER", true},
		{"USER_CREDENTIALS", true},
		{"DB_HOST", false},
		{"PORT", false},
		{"APP_ENV", false},
		{"LOG_LEVEL", false},
		// case insensitivity
		{"db_password", true},
		{"github_token", true},
	}
	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			got := IsSensitive(tc.key)
			if got != tc.wantSensitive {
				t.Errorf("IsSensitive(%q) = %v, want %v", tc.key, got, tc.wantSensitive)
			}
		})
	}
}

func TestRedactSensitiveEntries(t *testing.T) {
	input := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASSWORD", Value: "supersecret"},
		{Key: "PORT", Value: "5432"},
		{Key: "API_KEY", Value: "abc123"},
	}
	got := Redact(input)
	if len(got) != len(input) {
		t.Fatalf("expected %d entries, got %d", len(input), len(got))
	}
	if got[0].Value != "localhost" {
		t.Errorf("DB_HOST should not be redacted, got %q", got[0].Value)
	}
	if got[1].Value != RedactedPlaceholder {
		t.Errorf("DB_PASSWORD should be redacted, got %q", got[1].Value)
	}
	if got[2].Value != "5432" {
		t.Errorf("PORT should not be redacted, got %q", got[2].Value)
	}
	if got[3].Value != RedactedPlaceholder {
		t.Errorf("API_KEY should be redacted, got %q", got[3].Value)
	}
}

func TestRedactDoesNotMutateOriginal(t *testing.T) {
	input := []Entry{
		{Key: "API_KEY", Value: "original"},
	}
	_ = Redact(input)
	if input[0].Value != "original" {
		t.Errorf("Redact mutated original slice: got %q", input[0].Value)
	}
}

func TestRedactMap(t *testing.T) {
	m := map[string]string{
		"APP_NAME":    "envsync",
		"DB_PASSWORD": "hunter2",
		"GITHUB_TOKEN": "ghp_xxx",
	}
	got := RedactMap(m)
	if got["APP_NAME"] != "envsync" {
		t.Errorf("APP_NAME should not be redacted")
	}
	if got["DB_PASSWORD"] != RedactedPlaceholder {
		t.Errorf("DB_PASSWORD should be redacted, got %q", got["DB_PASSWORD"])
	}
	if got["GITHUB_TOKEN"] != RedactedPlaceholder {
		t.Errorf("GITHUB_TOKEN should be redacted, got %q", got["GITHUB_TOKEN"])
	}
	// original map unchanged
	if m["DB_PASSWORD"] != "hunter2" {
		t.Errorf("RedactMap mutated original map")
	}
}
