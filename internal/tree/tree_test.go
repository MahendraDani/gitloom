package tree_test

import (
	"os"
	"path/filepath"
	"strings"
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

func TestWriteTree_IgnoresGitloomDir(t *testing.T) {
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	filePath := filepath.Join(tempDir, "hello.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(filePath, content, repo.FilePerm); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	filePathWithinGitloom := filepath.Join(tempDir, repo.RepoDirName, "config")
	fileContent := []byte("config file\n")
	if err := os.WriteFile(filePathWithinGitloom, fileContent, repo.FilePerm); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	treeHash, err := tree.WriteTree(tempDir, r)
	if err != nil {
		t.Fatalf("WriteTree failed: %v", err)
	}

	objPath := filepath.Join(r.Path, repo.ObjectsDir, treeHash[:2], treeHash[2:])
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Fatalf("expected tree object to exist, but not found at %s", objPath)
	}

	output, err := object.CatFile(r, treeHash, "p")
	if err != nil {
		t.Fatalf("CatFile returned error: %v", err)
	}

	if !strings.Contains(output, "hello.txt") {
		t.Errorf("expected tree to contain 'hello.txt', got:\n%s", output)
	}

	if strings.Contains(output, ".gitloom") || strings.Contains(output, "config") {
		t.Errorf("tree contains entries from .gitloom directory, got:\n%s", output)
	}
}

func TestWriteTree_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	files := map[string]string{
		"a.txt": "alpha\n",
		"b.txt": "bravo\n",
	}

	for name, content := range files {
		path := filepath.Join(tempDir, name)
		if err := os.WriteFile(path, []byte(content), repo.FilePerm); err != nil {
			t.Fatalf("failed to write file %s: %v", name, err)
		}
	}

	treeHash, err := tree.WriteTree(tempDir, r)
	if err != nil {
		t.Fatalf("WriteTree returned error: %v", err)
	}

	objPath := filepath.Join(r.Path, repo.ObjectsDir, treeHash[:2], treeHash[2:])
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Fatalf("expected tree object to exist, but not found at %s", objPath)
	}

	objType, err := object.CatFile(r, treeHash, "t")
	if err != nil {
		t.Fatalf("CatFile -t returned error: %v", err)
	}
	if objType != "tree" {
		t.Fatalf("expected object type 'tree', got '%s'", objType)
	}

	output, err := object.CatFile(r, treeHash, "p")
	if err != nil {
		t.Fatalf("CatFile -p returned error: %v", err)
	}

	for name := range files {
		if !strings.Contains(output, name) {
			t.Errorf("expected tree to contain file '%s', got:\n%s", name, output)
		}
	}
}

func TestWriteTree_EmptyDirectory(t *testing.T) {
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("failed to read temp dir: %v", err)
	}
	for _, e := range entries {
		if e.Name() != repo.RepoDirName {
			t.Fatalf("expected only %s directory, found: %s", repo.RepoDirName, e.Name())
		}
	}

	treeHash, err := tree.WriteTree(tempDir, r)
	if err != nil {
		t.Fatalf("WriteTree returned error: %v", err)
	}

	objPath := filepath.Join(r.Path, repo.ObjectsDir, treeHash[:2], treeHash[2:])
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Fatalf("expected tree object to exist, but not found at %s", objPath)
	}

	objType, err := object.CatFile(r, treeHash, "t")
	if err != nil {
		t.Fatalf("CatFile -t returned error: %v", err)
	}
	if objType != "tree" {
		t.Fatalf("expected object type 'tree', got '%s'", objType)
	}

	output, err := object.CatFile(r, treeHash, "p")
	if err != nil {
		t.Fatalf("CatFile -p returned error: %v", err)
	}
	if output != "" {
		t.Fatalf("expected empty tree content, got:\n%s", output)
	}
}
