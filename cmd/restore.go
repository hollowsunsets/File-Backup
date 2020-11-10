package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use: "restore",
	Short: "Restore a file directory",
	Long: "Longer description here",
	Args: func(cmd *cobra.Command, args []string) error {
		if dirNameArg != "" && bucketNameArg != "" {
			return nil
		}
		if len(args) < 2 {
			return errors.New("directory name and bucket name required")
		}
		return nil
	},
	Run: func (cmd *cobra.Command, args[] string) {
		var directoryName string
		var bucketName string
		if dirNameArg != "" && bucketNameArg != "" {
			directoryName, bucketName = dirNameArg, bucketNameArg
		} else {
			directoryName, bucketName = args[0], args[1]
		}
		downloadBucket(directoryName, bucketName)
	},
}

func downloadBucket(directoryName string, bucketName string) {

}

func restoreDirectory(directoryName string, bucketName string) {

}