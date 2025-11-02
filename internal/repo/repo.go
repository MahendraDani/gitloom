package repo

import (
	"errors"
	"os"
	"path/filepath"
)

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

type Repo struct {
	Path string
}

func NewRepo(path string) *Repo {
	return &Repo{Path: path}
}

func FindRepo(startPath string) (*Repo, error) {
	path, err := filepath.Abs(startPath)
	if err != nil {
		return nil, err
	}
	for {
		repoPath := filepath.Join(path, RepoDirName)
		if _, err := os.Stat(repoPath); err == nil {
			return NewRepo(repoPath), nil
		}

		parent := filepath.Dir(path)
		if parent == path {
			return nil, errors.New("gitloom repository not found")
		}
		path = parent
	}
}

func (r *Repo) Init() error {
	path, err := filepath.Abs(r.Path)
	if err != nil {
		return err
	}

	repoPath := filepath.Join(path, RepoDirName)

	// check if repo already exists
	if _, err := os.Stat(repoPath); err == nil {
		return errors.New("repository already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	// create refs/head dir
	if err := os.MkdirAll(filepath.Join(repoPath, HeadsDir), DirPerm); err != nil {
		return err
	}

	// create refs/objects dir
	if err := os.MkdirAll(filepath.Join(repoPath, ObjectsDir), DirPerm); err != nil {
		return err
	}

	// create HEAD file
	headContent := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(filepath.Join(repoPath, HeadFile), headContent, FilePerm); err != nil {
		return err
	}

	r.Path = repoPath
	return nil
}

// TODO: implement them
// func (r *Repo) WriteObject(object internal.Object) (string, error)
// func (r *Repo) ReadObject(hash string) (internal.Object, error)
