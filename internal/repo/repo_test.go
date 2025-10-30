package repo_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MahendraDani/gitloom.git/internal/repo"
)

func TestInitRepository(t *testing.T) {
	tempDir := t.TempDir()

	// Run the function under test
	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("InitRepository failed: %v", err)
	}

	// Check .gitloom directory exists
	gitloomPath := filepath.Join(tempDir, repo.RepoDirName)
	if _, err := os.Stat(gitloomPath); os.IsNotExist(err) {
		t.Fatalf(".gitloom directory was not created")
	}

	// Check subdirectories
	refsPath := filepath.Join(gitloomPath, repo.RefsDir)
	if _, err := os.Stat(refsPath); os.IsNotExist(err) {
		t.Fatalf("refs directory was not created")
	}

	headsPath := filepath.Join(gitloomPath, repo.HeadsDir)
	if _, err := os.Stat(headsPath); os.IsNotExist(err) {
		t.Fatalf("refs/heads directory was not created")
	}

	objectsPath := filepath.Join(gitloomPath, repo.HeadsDir)
	if _, err := os.Stat(objectsPath); os.IsNotExist(err) {
		t.Fatalf("refs/objects directory was not created")
	}

	// Check HEAD file and content
	headFile := filepath.Join(gitloomPath, repo.HeadFile)
	data, err := os.ReadFile(headFile)
	if err != nil {
		t.Fatalf("failed to read HEAD file: %v", err)
	}

	expected := "ref: refs/heads/main\n"
	if string(data) != expected {
		t.Fatalf("unexpected HEAD file content:\n got: %q\nwant: %q", string(data), expected)
	}

	// Check return value correctness
	if r == nil {
		t.Fatalf("expected non-nil Repository")
	}
	if r.Path != gitloomPath {
		t.Fatalf("expected Repository.Path=%q, got %q", gitloomPath, r.Path)
	}
}

/*
1. Uses a temporary directory (t.TempDir()).
2. Manually creates a .gitloom/ directory inside it.
3. Calls InitRepository() on that same temp directory.
4. Verifies:
 1. The function returns an error.
 2. The error message matches "repository already exists"
*/
func TestInitRepositoryAlreadyExists(t *testing.T) {
	tempDir := t.TempDir()

	gitloomPath := filepath.Join(tempDir, repo.RepoDirName)
	if err := os.MkdirAll(gitloomPath, 0755); err != nil {
		t.Fatalf("failed to create fake .gitloom directory: %v", err)
	}

	_, err := repo.InitRepository(tempDir)
	if err == nil {
		t.Fatalf("expected error for existing repository, got nil")
	}
	expectedErr := "repository already exists"
	if err.Error() != expectedErr {
		t.Fatalf("unexpected error message:\n got: %q\nwant: %q", err.Error(), expectedErr)
	}
}

func TestInitRepositoryInvalidPath(t *testing.T) {
	// Provide an invalid path — e.g. a directory we can't create inside.
	// On Unix, /root usually requires root privileges.
	invalidPath := "/root/gitloom-test-invalid"

	_, err := repo.InitRepository(invalidPath)
	if err == nil {
		t.Fatalf("expected error for invalid path, got nil")
	}

	if err.Error() == "repository already exists" {
		t.Fatalf("unexpected 'repository already exists' error for invalid path")
	}
}
