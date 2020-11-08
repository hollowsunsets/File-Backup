package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)


func init() {
	rootCmd.AddCommand()
}

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "Backup a file directory",
	Long: "Longer description here",
	Run: func (cmd *cobra.Command, args[] string) {
		// directoryName, bucketName := args[0], args[1]
	},
}

/*
Find some directory
Find some bucket

Upload that directory to a bucket with that name


 */

func backup(directoryName string, bucketName string) {
	// Walks the given file tree in lexical order
	err := filepath.Walk(directoryName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			return nil
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}

