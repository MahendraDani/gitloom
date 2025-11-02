package cmd

import (
	"fmt"
	"log"

	"github.com/MahendraDani/gitloom.git/internal/object"
	"github.com/MahendraDani/gitloom.git/internal/repo"
	"github.com/spf13/cobra"
)

var (
	printFlag bool
	sizeFlag  bool
	typeFlag  bool
)

var catFileCmd = &cobra.Command{
	Use:   "cat-file <flag> <hash>",
	Short: "Provide content or information about objects in the gitloom datastore",
	Long: `gitloom cat-file displays the type, size, or contents of an object 
in the .gitloom/objects directory.

Usage:
  gitloom cat-file -p <hash>   # print object contents
  gitloom cat-file -s <hash>   # print object size
  gitloom cat-file -t <hash>   # print object type`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hash := args[0]

		r, err := repo.FindRepo(".")
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		flag := ""
		switch {
		case printFlag:
			flag = "p"
		case sizeFlag:
			flag = "s"
		case typeFlag:
			flag = "t"
		default:
			flag = "p" // default behavior
		}

		output, err := object.CatFile(r, hash, flag)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(catFileCmd)
	catFileCmd.Flags().BoolVarP(&printFlag, "print", "p", false, "Print object contents")
	catFileCmd.Flags().BoolVarP(&sizeFlag, "size", "s", false, "Print object size")
	catFileCmd.Flags().BoolVarP(&typeFlag, "type", "t", false, "Print object type")
}
