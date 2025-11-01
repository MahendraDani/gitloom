package object

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MahendraDani/gitloom.git/internal/repo"
)

func HashObject(filePath string, r *repo.Repository, write bool) (string, error) {
	if r == nil {
		return "", errors.New("gitloom repository not found. First initialize gitloom repository")
	}

	data, err := readFileBuffered(filePath)
	if err != nil {
		return "", err
	}

	blob := createBlob(data)
	hash := computeSHA1(blob)

	if write {
		if err := writeObject(blob, hash, r); err != nil {
			return "", err
		}
	}

	return hash, nil
}

func readFileBuffered(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var buf bytes.Buffer
	reader := bufio.NewReader(f)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func createBlob(data []byte) []byte {
	header := []byte(fmt.Sprintf("blob %d\x00", len(data)))
	return append(header, data...)
}

func computeSHA1(blob []byte) string {
	h := sha1.New()
	h.Write(blob)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func writeObject(blob []byte, hash string, r *repo.Repository) error {
	objDir := filepath.Join(r.Path, repo.ObjectsDir, hash[:2])
	objPath := filepath.Join(objDir, hash[2:])

	// Ensure directory exists
	if err := os.MkdirAll(objDir, repo.DirPerm); err != nil {
		return err
	}

	// Use helper to compress & write
	return writeZlibFile(objPath, blob)
}

func writeZlibFile(filePath string, data []byte) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zlib.NewWriter(f)
	defer zw.Close()

	if _, err := zw.Write(data); err != nil {
		return err
	}

	return nil
}

func CatFile(r *repo.Repository, hash string, flag string) (string, error) {
	if r == nil {
		return "", errors.New("gitloom repository not found")
	}

	if len(hash) < 2 {
		return "", errors.New("invalid object hash")
	}

	objPath := filepath.Join(r.Path, repo.ObjectsDir, hash[:2], hash[2:])
	data, err := DecompressObject(objPath)
	if err != nil {
		return "", err
	}

	nullIdx := bytes.IndexByte(data, 0)
	if nullIdx == -1 {
		return "", errors.New("invalid object format (missing header)")
	}

	header := string(data[:nullIdx])
	content := string(data[nullIdx+1:])

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid object header format")
	}

	objType := parts[0]

	switch flag {
	case "p":
		if objType != "blob" {
			return "", fmt.Errorf("cat-file -p only supports blob objects, got %s", objType)
		}
		return content, nil

	case "s":
		if objType != "blob" {
			return "", fmt.Errorf("cat-file -s only supports blob objects, got %s", objType)
		}
		return fmt.Sprintf("%d", len(content)), nil
	case "t":
		return objType, nil
	default:
		return "", fmt.Errorf("unsupported flag: %s", flag)
	}
}

func DecompressObject(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open object file: %w", err)
	}
	defer f.Close()

	zr, err := zlib.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer zr.Close()

	br := bufio.NewReader(zr)
	var buf bytes.Buffer

	for {
		chunk := make([]byte, 4096)
		n, err := br.Read(chunk)
		if n > 0 {
			buf.Write(chunk[:n])
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed while reading compressed data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func HashRawObject(data []byte, objType string, r *repo.Repository, write bool) (string, error) {
	// Build header: "<type> <size>\0"
	header := fmt.Sprintf("%s %d\x00", objType, len(data))
	store := append([]byte(header), data...)

	// Compute SHA-1 hash of the full content
	hash := sha1.Sum(store)
	hashHex := hex.EncodeToString(hash[:])

	// If write == false, just return hash
	if !write {
		return hashHex, nil
	}

	// Construct object path (.gitloom/objects/xx/yyyy...)
	objDir := filepath.Join(r.Path, repo.ObjectsDir, hashHex[:2])
	objPath := filepath.Join(objDir, hashHex[2:])

	// Avoid rewriting existing objects
	if _, err := os.Stat(objPath); err == nil {
		return hashHex, nil
	}

	// Ensure object subdirectory exists
	if err := os.MkdirAll(objDir, repo.DirPerm); err != nil {
		return "", err
	}

	// Compress and write object data
	writeZlibFile(objPath, store)
	return hashHex, nil
}
