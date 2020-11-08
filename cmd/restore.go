package cmd

import "github.com/spf13/cobra"

func restore(bucketName string, directoryName string) {

}

var restoreCmd = &cobra.Command{
	Use: "restore",
	Short: "Restore a file directory",
	Long: "Longer description here",
	Run: func (cmd *cobra.Command, args[] string) {
		// directoryName, bucketName := args[0], args[1]
	},
}
