package object_test

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
)

func TestHashObjectAndWrite(t *testing.T) {
	// Create a temporary directory for the repo
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	fmt.Println(r)

	// Create a test file
	filePath := filepath.Join(tempDir, "hello.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	writeFlag := true
	// Call HashObject
	hash, err := object.HashObject(filePath, r, writeFlag)
	if err != nil {
		t.Fatalf("HashObject returned error: %v", err)
	}

	// Check that the blob file exists in .gitloom/objects
	objPath := filepath.Join(r.Path, repo.ObjectsDir, hash[:2], hash[2:])
	if _, err := os.Stat(objPath); err != nil {
		t.Fatalf("expected object file at %s, but got error: %v", objPath, err)
	}

	// Optional: check that the content is correct (zlib decompressed)
	objFile, err := os.Open(objPath)
	if err != nil {
		t.Fatalf("failed to open object file: %v", err)
	}
	defer objFile.Close()

	zr, err := zlib.NewReader(objFile)
	if err != nil {
		t.Fatalf("failed to create zlib reader: %v", err)
	}
	defer zr.Close()

	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, zr); err != nil {
		t.Fatalf("failed to decompress object: %v", err)
	}

	expectedHeader := []byte(fmt.Sprintf("blob %d\x00", len(content)))
	expectedBlob := append(expectedHeader, content...)

	if !bytes.Equal(decompressed.Bytes(), expectedBlob) {
		t.Fatalf("object content mismatch:\n got:  %q\n want: %q", decompressed.Bytes(), expectedBlob)
	}
}

func TestHashObject(t *testing.T) {
	// Create a temporary directory for the repo
	tempDir := t.TempDir()

	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	fmt.Println(r)

	// Create a test file
	filePath := filepath.Join(tempDir, "hello.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	writeFlag := false
	// Call HashObject
	hash, err := object.HashObject(filePath, r, writeFlag)
	if err != nil {
		t.Fatalf("HashObject returned error: %v", err)
	}

	objPath := filepath.Join(r.Path, repo.ObjectsDir, hash[:2], hash[2:])
	if _, err := os.Stat(objPath); !os.IsNotExist(err) {
		t.Fatalf("expected object NOT to exists, but found at %s", objPath)
	}
}

func TestCatFilePrint(t *testing.T) {
	// Create a temporary directory for the repo
	tempDir := t.TempDir()

	// setup a gitloom repo
	r, err := repo.InitRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to init repository: %v", err)
	}

	// create a temp file
	filePath := filepath.Join(tempDir, "hello.txt")
	content := []byte("hello world\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// write the file into `.gitloom/objects/*`
	writeFlag := true
	hash, err := object.HashObject(filePath, r, writeFlag)
	if err != nil {
		t.Fatalf("HashObject returned error: %v", err)
	}

	flag := "p"
	output, err := object.CatFile(r, hash, flag)
	if err != nil {
		t.Fatalf("CatFile returned error: %v", err)
	}

	expected := string(content)
	if output != expected {
		t.Fatalf("expected content %q, got %q", expected, output)
	}
}
