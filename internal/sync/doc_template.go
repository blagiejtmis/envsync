// Package sync provides push/pull synchronisation of .env files.
//
// Template checking
//
// The template_check.go file adds the ability to compare a live .env file
// against a committed .env.template, detecting:
//
//   - Missing keys  – declared in the template but absent or blank in the
//     local environment file. These must be filled before a push is accepted.
//
//   - Extra keys    – present in the local file but not declared in the
//     template. These may indicate stale or undocumented configuration.
//
// Typical workflow:
//
//  1. A developer runs `envsync template generate` to produce .env.template
//     from the current .env (values are stripped).
//  2. The template is committed to version control.
//  3. On each push/pull, CheckAgainstTemplate is called automatically and
//     the result is reported via FprintCheckResult.
package sync
