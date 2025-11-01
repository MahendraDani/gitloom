package tree

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
)

// WriteTree creates a tree object representing the current state
// of the working directory (ignoring .gitloom and subdirectories for now).
func WriteTree(dir string, r *repo.Repository) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var treeBuf bytes.Buffer

	for _, entry := range entries {
		name := entry.Name()

		// Ignore the .gitloom directory
		if name == repo.RepoDirName {
			continue
		}

		// Handle only regular files for now
		if entry.Type().IsRegular() {
			filePath := filepath.Join(dir, name)

			// Create a blob object for this file
			blobHash, err := object.HashObject(filePath, r, true)
			if err != nil {
				return "", err
			}

			// Convert hex string hash → raw bytes
			hashBytes, err := hex.DecodeString(blobHash)
			if err != nil {
				return "", err
			}

			// Write tree entry: "<mode> <filename>\0<binary hash>"
			treeBuf.WriteString("100644 " + name + "\x00")
			treeBuf.Write(hashBytes)
		}
	}

	// Finally, write this tree object to the `.gitloom/objects/*`
	treeHash, err := object.HashRawObject(treeBuf.Bytes(), "tree", r, true)
	if err != nil {
		return "", err
	}

	return treeHash, nil
}
