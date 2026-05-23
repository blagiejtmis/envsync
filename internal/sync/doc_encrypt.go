// Package sync provides functionality for pushing and pulling .env files
// to and from a shared secret store.
//
// # Field-Level Encryption
//
// Field-level encryption allows individual sensitive values to be encrypted
// before they are written to the secret store, providing an additional layer
// of protection beyond the store's own access controls.
//
// Encrypted values are stored with the "enc:" prefix so they can be
// identified and decrypted transparently on pull.
//
//	applyEncryptPush(entries, EncryptPushOptions{
//	    Enabled:    true,
//	    Passphrase: "my-team-secret",
//	})
//
// By default, any key deemed sensitive by env.IsSensitive is encrypted.
// Specific keys can also be listed explicitly via EncryptFieldOptions.Keys.
//
// The same passphrase must be supplied on pull to decrypt the values:
//
//	applyDecryptPull(entries, "my-team-secret")
package sync
