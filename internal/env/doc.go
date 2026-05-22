// Package env provides utilities for working with .env files throughout the
// envsync toolchain.
//
// Capabilities:
//
//   - Parse / Serialize: read and write KEY=VALUE lines, honouring comments and
//     blank lines.
//   - LoadFile / WriteFile: thin filesystem helpers built on top of Parse and
//     Serialize.
//   - Diff / Merge: compare two sets of entries and merge them with configurable
//     conflict resolution strategies.
//   - Validate: enforce key-naming rules and detect blank values.
//   - Redact: mask sensitive values before logging or displaying output.
//   - Export: render entries in raw, export, or Docker --env-file formats.
//   - Snapshot: capture a point-in-time checksum of an env file for integrity
//     verification.
//   - Watch: poll an env file for changes and emit ChangeEvents over a channel,
//     enabling live-reload workflows during development.
package env
