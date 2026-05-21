// Package sync provides the core push/pull synchronisation logic for envsync.
//
// It ties together the crypto, env, and store packages to allow team members
// to securely share .env files via a shared secret store backend.
//
// # Usage
//
//	s, _ := store.New("/path/to/store.json")
//	syncer := sync.New(s, "shared-passphrase")
//
//	// Encrypt and upload
//	syncer.Push("myapp-prod", ".env")
//
//	// Download and decrypt
//	syncer.Pull("myapp-prod", ".env")
//
// The payload stored in the backend is: [16-byte salt][AES-GCM ciphertext].
// The salt is randomly generated on each Push, ensuring that the same plaintext
// produces a different ciphertext every time.
package sync
