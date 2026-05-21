package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envsync/internal/store"
	"github.com/yourorg/envsync/internal/sync"
)

const defaultStorePath = ".envsync_store.json"

func main() {
	if len(os.Args) < 4 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	key := os.Args[2]
	filePath := os.Args[3]

	passphrase := os.Getenv("ENVSYNC_PASSPHRASE")
	if passphrase == "" {
		fmt.Fprintln(os.Stderr, "error: ENVSYNC_PASSPHRASE environment variable is not set")
		os.Exit(1)
	}

	storePath := os.Getenv("ENVSYNC_STORE")
	if storePath == "" {
		storePath = defaultStorePath
	}

	s, err := store.New(storePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: open store: %v\n", err)
		os.Exit(1)
	}

	syncer := sync.New(s, passphrase)

	switch command {
	case "push":
		if err := syncer.Push(key, filePath); err != nil {
			fmt.Fprintf(os.Stderr, "error: push: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("pushed %s -> store[%s]\n", filePath, key)

	case "pull":
		if err := syncer.Pull(key, filePath); err != nil {
			fmt.Fprintf(os.Stderr, "error: pull: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("pulled store[%s] -> %s\n", key, filePath)

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage:
  envsync push <key> <file>   Encrypt and push .env file to store
  envsync pull <key> <file>   Pull and decrypt .env file from store

Environment variables:
  ENVSYNC_PASSPHRASE   Shared secret passphrase (required)
  ENVSYNC_STORE        Path to store file (default: %s)
`, defaultStorePath)
}
