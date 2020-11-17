package cmd

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hollowsunsets/fcopy/utils"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"time"
)


func init() {
	rootCmd.AddCommand()
}

var backupCmd = &cobra.Command{
	Use: "backup",
	Short: "Backup a file directory",
	Long: `Backup a local file directory to a AWS S3 bucket.
Ex: fcopy backup <directory_name> <bucket_name>
    fcopy backup --bucket <bucket_name> --directory <directory_name>
If no S3 bucket named <bucket_name> exists, a new S3 bucket named <bucket_name> will be created.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if dirNameArg != "" && bucketNameArg != "" {
			return nil
		}
		if len(args) < 2 {
			cmd.Help()
			os.Exit(0)
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

		// If the given directory doesn't exist
		if _, err := os.Stat(directoryName); os.IsNotExist(err) {
			fmt.Printf("given directory %s does not exist\n", directoryName)
			os.Exit(1)
		}

		err := backupDirectory(directoryName, bucketName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func backupDirectory(directoryName string, bucketName string) error {
	// Walks the given file tree in lexical order
	err := filepath.Walk(directoryName,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err = uploadFile(path, bucketName, info.ModTime())
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	return err
}

func uploadFile(fileName string, bucketName string, modTime time.Time) error {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %s: \"%s\"", fileName, err)
	}
	defer f.Close()

	if ok := utils.BucketExists(bucketName); !ok {
		err := utils.CreateBucket(bucketName)
		if err != nil {
			return fmt.Errorf("error creating bucket %s: \"%s\"", bucketName, err)
		}
	}
	metadata, err := utils.GetObjectMetadata(bucketName, fileName)
	if err == nil {
		if !modTime.After(*metadata.LastModified) {
			fmt.Printf("file %s has not been updated, no need to backup\n", fileName)
			return nil
		}
	}

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName),
		Body:   f,
	}
	result, err := uploader.Upload(input)
	if err != nil {
		return fmt.Errorf("error uploading file to bucket %s: \"%s\"", bucketName, err)
	}
	fmt.Printf("file uploaded to %s\n", result.Location)
	return nil
}
