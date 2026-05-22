// Package version exposes build-time version metadata for the envsync binary.
//
// Version fields (Major, Minor, Patch, PreRelease, Commit, BuildDate) are
// intended to be overridden at link time using -ldflags, for example:
//
//	-ldflags "-X github.com/yourorg/envsync/internal/version.Major=1 \
//	          -X github.com/yourorg/envsync/internal/version.Minor=2 \
//	          -X github.com/yourorg/envsync/internal/version.Patch=3 \
//	          -X github.com/yourorg/envsync/internal/version.PreRelease=beta.1 \
//	          -X github.com/yourorg/envsync/internal/version.Commit=$(git rev-parse --short HEAD) \
//	          -X github.com/yourorg/envsync/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)"
//
// When no flags are provided the package falls back to safe defaults so that
// development builds remain functional.
//
// Typical usage:
//
//	fmt.Println(version.String())       // "1.2.3-beta.1"
//	fmt.Println(version.Full())         // "1.2.3-beta.1 (commit=abc1234, built=2024-01-15T10:00:00Z)"
package version
