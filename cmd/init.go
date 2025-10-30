package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/spf13/cobra"
)

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

		if _, err := repo.InitRepository(absPath); err != nil {
			fmt.Println("Error initializing repository:", err)
			return
		}

		fmt.Println("Initialized empty gitloom repository at", path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
