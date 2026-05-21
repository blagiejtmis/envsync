// Package audit provides append-only audit logging for envsync operations.
//
// Each push, pull, or delete action on a secret key can be recorded as a
// JSON-lines entry in a local log file. Entries include a UTC timestamp,
// the operation kind, the affected key name, an optional username, and an
// optional free-form message.
//
// Example usage:
//
//	logger := audit.NewLogger("/var/log/envsync/audit.log")
//	logger.Record(audit.EventPush, "DB_PASSWORD", "alice", "initial push")
//
// The log file is created automatically if it does not exist and is opened
// in append mode so concurrent writes from multiple processes are safe on
// POSIX systems.
package audit
