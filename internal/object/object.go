package object

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MahendraDani/gitloom.git/internal/repo"
)

func HashObject(filePath string, r *repo.Repository) (string, error) {
	if r == nil {
		return "", errors.New("gitloom repository not found. First initialize gitloom repository")
	}

	data, err := readFileBuffered(filePath)
	if err != nil {
		return "", err
	}

	blob := createBlob(data)
	hash := computeSHA1(blob)

	if err := writeObject(blob, hash, r); err != nil {
		return "", err
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
