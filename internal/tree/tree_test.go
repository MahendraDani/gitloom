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

	r := repo.NewRepo(tempDir)
	if err := r.Init(); err != nil {
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

	r := repo.NewRepo(tempDir)
	if err := r.Init(); err != nil {
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

	r := repo.NewRepo(tempDir)
	if err := r.Init(); err != nil {
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

	r := repo.NewRepo(tempDir)
	if err := r.Init(); err != nil {
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

func TestWriteTree_WithSubdirectories(t *testing.T) {
	tempDir := t.TempDir()

	r := repo.NewRepo(tempDir)
	if err := r.Init(); err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}
	// Create files
	rootFile := filepath.Join(tempDir, "file1.txt")
	if err := os.WriteFile(rootFile, []byte("root content\n"), repo.FilePerm); err != nil {
		t.Fatalf("failed to write root file: %v", err)
	}

	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, repo.DirPerm); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	subFile1 := filepath.Join(subDir, "file2.txt")
	subFile2 := filepath.Join(subDir, "file3.txt")

	if err := os.WriteFile(subFile1, []byte("sub content 1\n"), repo.FilePerm); err != nil {
		t.Fatalf("failed to write sub file1: %v", err)
	}
	if err := os.WriteFile(subFile2, []byte("sub content 2\n"), repo.FilePerm); err != nil {
		t.Fatalf("failed to write sub file2: %v", err)
	}

	// Call WriteTree
	rootTreeHash, err := tree.WriteTree(tempDir, r)
	if err != nil {
		t.Fatalf("WriteTree returned error: %v", err)
	}

	// Verify root tree object exists
	rootObjPath := filepath.Join(r.Path, repo.ObjectsDir, rootTreeHash[:2], rootTreeHash[2:])
	if _, err := os.Stat(rootObjPath); os.IsNotExist(err) {
		t.Fatalf("expected root tree object to exist, not found at %s", rootObjPath)
	}

	// Verify root tree type
	rootType, err := object.CatFile(r, rootTreeHash, "t")
	if err != nil {
		t.Fatalf("CatFile -t for root tree returned error: %v", err)
	}
	if rootType != "tree" {
		t.Fatalf("expected root tree object type 'tree', got '%s'", rootType)
	}

	// Pretty-print root tree
	rootTree, err := object.CatFile(r, rootTreeHash, "p")
	if err != nil {
		t.Fatalf("CatFile -p for root tree returned error: %v", err)
	}

	if !strings.Contains(rootTree, "file1.txt") {
		t.Fatalf("expected root tree to contain 'file1.txt', got:\n%s", rootTree)
	}
	if !strings.Contains(rootTree, "subdir") {
		t.Fatalf("expected root tree to contain 'subdir', got:\n%s", rootTree)
	}

	// Extract the hash of subdir tree (optional deeper validation)
	lines := strings.Split(rootTree, "\n")
	var subdirHash string
	for _, line := range lines {
		if strings.Contains(line, "subdir") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				subdirHash = fields[2]
				break
			}
		}
	}

	if subdirHash == "" {
		t.Fatalf("failed to extract subdir hash from root tree output")
	}

	// Verify sub-tree type and content
	subType, err := object.CatFile(r, subdirHash, "t")
	if err != nil {
		t.Fatalf("CatFile -t for subdir returned error: %v", err)
	}
	if subType != "tree" {
		t.Fatalf("expected subdir object type 'tree', got '%s'", subType)
	}

	subTree, err := object.CatFile(r, subdirHash, "p")
	if err != nil {
		t.Fatalf("CatFile -p for subdir returned error: %v", err)
	}

	if !strings.Contains(subTree, "file2.txt") || !strings.Contains(subTree, "file3.txt") {
		t.Fatalf("expected subdir tree to contain file2.txt and file3.txt, got:\n%s", subTree)
	}
}
