package repo

type Repository struct {
	Path string
}

const (
	RepoDirName = ".gitloom"
	HeadFile    = "HEAD"
	RefsDir     = "refs"
	HeadsDir    = "refs/heads"
	MainBranch  = "main"
)

/*
- check if .gitloom dir already exists => return error
- if not, then
  - create a .gitloom directory within the provided path
  - create a HEAD file within .gitloom/
  - the file should contain one-line: ref: refs/heads/main
  - create a refs directory => .gitloom/refs
  - create a heads direcotyr => .gitloom/refs/heads
*/
func InitRepository(path string) (*Repository, error) {
	return &Repository{
		Path: "write-code-here",
	}, nil
}
