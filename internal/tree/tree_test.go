package tree_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/MahendraDani/gitloom.git/internal/tree"
)

func TestWriteTree_SingleFile(t *testing.T) {
	// Create a temporary directory for the repo
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	// Create a test file in the root
	filePath := filepath.Join(tempDir, "hello.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(filePath, content, repo.FilePerm); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Call WriteTree
	treeHash, err := tree.WriteTree(tempDir, r)
	if err != nil {
		t.Fatalf("WriteTree returned error: %v", err)
	}

	// Verify that the corresponding tree object exists
	objPath := filepath.Join(r.Path, repo.ObjectsDir, treeHash[:2], treeHash[2:])
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Fatalf("expected tree object to exist, but not found at %s", objPath)
	}

	// Verify that the object is of type "tree" using CatFile
	objType, err := object.CatFile(r, treeHash, "t")
	if err != nil {
		t.Fatalf("CatFile returned error: %v", err)
	}

	if objType != "tree" {
		t.Fatalf("expected object type 'tree', got '%s'", objType)
	}
}
