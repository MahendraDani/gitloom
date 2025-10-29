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
