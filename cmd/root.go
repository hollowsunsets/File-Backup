package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var (
	directoryName string
	bucketName string
	rootCmd = &cobra.Command{
		Use: "fcopy",
		Short: "fcopy is a CLI for that manages file backups and restores in AWS S3.",
		Long: `fcopy is a CLI that can upload and retrieve backups of file directories 
to and from a bucket in Amazon's Simple Storage Service (AWS S3).`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Hello CLI")
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&directoryName, "directory", "", "Name of file directory")
	rootCmd.PersistentFlags().StringVar(&bucketName, "bucket", "", "Name of AWS S3 bucket")
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(restoreCmd)
}