package version_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envsync/internal/version"
)

func TestGetReturnsInfo(t *testing.T) {
	info := version.Get()
	if info.Major == "" {
		t.Error("expected non-empty Major")
	}
	if info.Minor == "" {
		t.Error("expected non-empty Minor")
	}
	if info.Patch == "" {
		t.Error("expected non-empty Patch")
	}
}

func TestInfoString(t *testing.T) {
	info := version.Info{
		Major:      "1",
		Minor:      "2",
		Patch:      "3",
		PreRelease: "beta",
		Commit:     "abc1234",
		BuildDate:  "2024-01-15",
	}

	got := info.String()
	if got != "1.2.3-beta" {
		t.Errorf("String() = %q, want %q", got, "1.2.3-beta")
	}
}

func TestInfoStringNoPreRelease(t *testing.T) {
	info := version.Info{
		Major: "2",
		Minor: "0",
		Patch: "0",
	}

	got := info.String()
	if got != "2.0.0" {
		t.Errorf("String() = %q, want %q", got, "2.0.0")
	}
	if strings.Contains(got, "-") {
		t.Errorf("String() should not contain '-' when PreRelease is empty, got %q", got)
	}
}

func TestInfoFull(t *testing.T) {
	info := version.Info{
		Major:      "0",
		Minor:      "1",
		Patch:      "0",
		PreRelease: "dev",
		Commit:     "deadbeef",
		BuildDate:  "2024-06-01",
	}

	got := info.Full()
	if !strings.Contains(got, "deadbeef") {
		t.Errorf("Full() missing commit hash, got %q", got)
	}
	if !strings.Contains(got, "2024-06-01") {
		t.Errorf("Full() missing build date, got %q", got)
	}
	if !strings.Contains(got, "0.1.0-dev") {
		t.Errorf("Full() missing version string, got %q", got)
	}
}

func TestDefaultVersionIsDevBuild(t *testing.T) {
	info := version.Get()
	if info.PreRelease != "dev" {
		t.Errorf("default PreRelease = %q, want %q", info.PreRelease, "dev")
	}
	if info.Commit != "unknown" {
		t.Errorf("default Commit = %q, want %q", info.Commit, "unknown")
	}
}
