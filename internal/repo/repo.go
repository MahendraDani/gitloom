package repo

import (
	"errors"
	"os"
	"path/filepath"
)

type Repository struct {
	Path string // store absolute path - path/to/project/.gitloom
}

const (
	RepoDirName = ".gitloom"
	HeadFile    = "HEAD"
	RefsDir     = "refs"
	HeadsDir    = "refs/heads"
	ObjectsDir  = "objects"
	MainBranch  = "main"

	DirPerm  = 0755
	FilePerm = 0644
)

/*
- check if .gitloom dir already exists => return error
- if not, then
  - create a .gitloom directory within the provided path
  - create a HEAD file within .gitloom/
  - the file should contain one-line: ref: refs/heads/main
  - create a refs directory => .gitloom/refs
  - create a objects directory => .gitloom/objects
  - create a heads direcotry => .gitloom/refs/heads
*/
func InitRepository(path string) (*Repository, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	repoPath := filepath.Join(path, RepoDirName)

	if _, err := os.Stat(repoPath); err == nil {
		return nil, errors.New("repository already exists")
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// create the res/heads dir
	if err := os.MkdirAll(filepath.Join(repoPath, HeadsDir), DirPerm); err != nil {
		return nil, err
	}

	// create .gitloom/objects dir
	if err := os.MkdirAll(filepath.Join(repoPath, ObjectsDir), DirPerm); err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(repoPath, HeadFile), []byte("ref: refs/heads/main\n"), FilePerm); err != nil {
		return nil, err
	}
	return &Repository{
		Path: repoPath,
	}, nil
}

/*
searches for a Gitloom repository starting from the given path and moving up the directory tree until it finds the repository root.

why?
so we need to get the repo struct, in different places in code, so we can just recreate it
by walking up the path and checking if we get the repo
*/
func FindRepository(startPath string) (*Repository, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}

	for {
		repoPath := filepath.Join(path, RepoDirName)
		if _, err := os.Stat(repoPath); err == nil {
			return &Repository{Path: repoPath}, nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return nil, errors.New("No gitloom repository found")
		}
		path = parent
	}
}
