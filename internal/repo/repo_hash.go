package repo

import (
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func HashObject(filePath string, r *Repository) (string, error) {
	if r == nil {
		return "", errors.New("gitloom repository not found. First initialize gitloom repository")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// For now, just return the length as string to test step
	header := fmt.Sprintf("blob %d\x00", len(data))
	blob := append([]byte(header), data...)
	h := sha1.New()
	h.Write(blob)
	hash := fmt.Sprintf("%x", h.Sum(nil))

	objDir := filepath.Join(r.Path, ObjectsDir, hash[:2])
	objPath := filepath.Join(objDir, hash[2:])

	if err := os.MkdirAll(objDir, DirPerm); err != nil {
		return "", err
	}

	f, err := os.Create(objPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	zw := zlib.NewWriter(f)
	defer zw.Close()

	if _, err := zw.Write(blob); err != nil {
		return "", err
	}
	return hash, nil
}
