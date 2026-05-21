// Package version provides build-time version information for envsync.
package version

import "fmt"

// These variables are set at build time via -ldflags.
var (
	// Major is the major version number.
	Major = "0"
	// Minor is the minor version number.
	Minor = "1"
	// Patch is the patch version number.
	Patch = "0"
	// PreRelease is an optional pre-release label (e.g. "alpha", "beta").
	PreRelease = "dev"
	// Commit is the short git commit hash injected at build time.
	Commit = "unknown"
	// BuildDate is the UTC build timestamp injected at build time.
	BuildDate = "unknown"
)

// Info holds structured version metadata.
type Info struct {
	Major      string
	Minor      string
	Patch      string
	PreRelease string
	Commit     string
	BuildDate  string
}

// Get returns the current version Info.
func Get() Info {
	return Info{
		Major:      Major,
		Minor:      Minor,
		Patch:      Patch,
		PreRelease: PreRelease,
		Commit:     Commit,
		BuildDate:  BuildDate,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	v := fmt.Sprintf("%s.%s.%s", i.Major, i.Minor, i.Patch)
	if i.PreRelease != "" {
		v += "-" + i.PreRelease
	}
	return v
}

// Full returns a verbose version string including commit and build date.
func (i Info) Full() string {
	return fmt.Sprintf("%s (commit=%s, built=%s)", i.String(), i.Commit, i.BuildDate)
}
