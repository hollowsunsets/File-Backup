package cmd

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if dirNameArg != "" && bucketNameArg != "" {
			return nil
		}
		if len(args) < 2 {
			return errors.New("directory name and bucket name required")
		}
		return nil
	},
	Run: func (cmd *cobra.Command, args []string) {
		var directoryName string
		var bucketName string
		if dirNameArg != "" && bucketNameArg != "" {
			directoryName, bucketName = dirNameArg, bucketNameArg
		} else {
			directoryName, bucketName = args[0], args[1]
		}
		findDirectory(directoryName, bucketName)
	},
}


func findDirectory(directoryName string, bucketName string) error {
	// Walks the given file tree in lexical order
	err := filepath.Walk(directoryName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			uploadFile(directoryName, bucketName)
			return nil
		},
	)
	return err
}

func uploadFile(fileName string, bucketName string) error {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", fileName)
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Body: f,
	})

	fmt.Printf("file uploaded to, %s\n\n", aws.StringValue(&result.Location))
	return nil
}