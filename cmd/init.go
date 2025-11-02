package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/spf13/cobra"
)

type ctxKey string

const repoKey ctxKey = "repo"

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new gitloom repository.",
	Long: `Initialize a new gitloom repository in a current working directory if no directory name is provided. For example:

gitloom init - creates a new gitloom repository within current working directory
gitloom init dir-name - creates a new gitloom repository within dir-name directory. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		var path string

		if len(args) > 0 {
			path = args[0]
		} else {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}
			path = cwd
		}

		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Println("Error resolving absolute path:", err)
			return
		}

		repo := repo.NewRepo(absPath)
		if err := repo.Init(); err != nil {
			fmt.Println("Error initializing repository:", err)
		}

		fmt.Println("Initialized empty gitloom repository at", repo.Path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
