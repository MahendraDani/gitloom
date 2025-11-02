package tree

import (
	"bytes"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
)

// WriteTree creates a tree object representing the current state
// of the working directory, recursively including subdirectories.
func WriteTree(dir string, r *repo.Repo) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var treeBuf bytes.Buffer

	// Sort entries to ensure deterministic order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		name := entry.Name()

		// Ignore the .gitloom directory
		if name == repo.RepoDirName {
			continue
		}

		fullPath := filepath.Join(dir, name)

		if entry.IsDir() {
			// Recursively write subdirectory tree
			subTreeHash, err := WriteTree(fullPath, r)
			if err != nil {
				return "", err
			}

			hashBytes, err := hex.DecodeString(subTreeHash)
			if err != nil {
				return "", err
			}

			// Mode 40000 for directories
			treeBuf.WriteString("40000 " + name + "\x00")
			treeBuf.Write(hashBytes)

		} else if entry.Type().IsRegular() {
			// Create blob object for the file
			blobHash, err := object.HashObject(fullPath, r, true)
			if err != nil {
				return "", err
			}

			hashBytes, err := hex.DecodeString(blobHash)
			if err != nil {
				return "", err
			}

			// Mode 100644 for normal files
			treeBuf.WriteString("100644 " + name + "\x00")
			treeBuf.Write(hashBytes)
		}
	}

	// Write the tree object to the .gitloom/objects directory
	treeHash, err := object.HashRawObject(treeBuf.Bytes(), "tree", r, true)
	if err != nil {
		return "", err
	}

	return treeHash, nil
}
