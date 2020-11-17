package cmd

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
		err := clearDirectory(directoryName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = downloadBucket(directoryName, bucketName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}



func clearDirectory(directoryName string) error {
	dir, err := os.Open(directoryName)
	if err != nil {
		return err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(directoryName, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func listBucketItems(bucketName string) (*s3.ListObjectsOutput, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket:  aws.String(bucketName),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return nil, err
	}
	return result, err

}

func downloadBucket(directoryName string, bucketName string) error {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)

	bucketObjects, err := listBucketItems(bucketName)
	if err != nil {
		return err
	}
	for _, object := range bucketObjects.Contents {
		filename := filepath.Join(directoryName, *object.Key)
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("failed to create file %q, %v", filename, err)
		}

		_, err = downloader.Download(f, &s3.GetObjectInput{
			Bucket: aws.String(bucketName),
			Key:	aws.String(filename),
		})
		if err != nil {
			return fmt.Errorf("failed to download file %q, %v", filename, err)
		}
	}
	return nil
}
