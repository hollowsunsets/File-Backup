package cmd

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/hollowsunsets/fcopy/utils"
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
		err := findDirectory(directoryName, bucketName)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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
			if !info.IsDir() {
				err = uploadFile(path, info, bucketName)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	return err
}

func uploadFile(fileName string, info os.FileInfo, bucketName string) error {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file %s: \"%s\"", fileName, err)
	}
	defer f.Close()

	if ok, err := backupWillNotChange(info, bucketName, fileName); ok && err != nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("error checking if backup and file are the same:\"%s\" ", err)
	}

	// If the MD5 digest is the same for both files,
	fileMD5, err := utils.FileMD5Checksum(f)
	input := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key : aws.String(fileName),
		Body: f,
		Metadata: map[string]*string{
			"fcopy-MD5": &fileMD5,
		},
	}
	result, err := uploader.Upload(input)

	if err != nil {
		return fmt.Errorf("error uploading file to bucket %s: \"%s\"", bucketName, err)
	}
	fmt.Printf("file uploaded to %s\n", result.Location)
	return nil
}


func backupWillNotChange(info os.FileInfo, bucketName string, fileName string) (bool, error) {
	metadata, err := utils.GetObjectMetadata(bucketName, fileName)
	if err != nil {
		return false, err
	}
	if info.ModTime().After(*metadata.LastModified) {
		return false, nil
	}

	object, err := utils.GetObject(bucketName, fileName)
	if err != nil {
		return true, err
	}

	f, err := os.Open(fileName)
	if err != nil {
		return true, err
	}

	// Check if the MD5 digest is the same.
	// If the MD5 digest is the same for both files,
	fileMD5, err := utils.FileMD5Checksum(f)
	if err != nil {
		return true, err
	}

	if ok, err := matchObjectMD5Sum(object, metadata, fileMD5); ok && err != nil {
		return true, nil
	}
	if err != nil {
		return true, err
	}
	return false, nil
}


func matchObjectMD5Sum(object *s3.GetObjectOutput, metadata *s3.HeadObjectOutput, fileMD5 string) (bool, error) {
	// Case: Native S3 entity tag exists, which may be an MD5 digest.
	// 		 If the object was uploaded in multiple parts, it will not have an MD5 digest.
	// 		 See: https://forums.aws.amazon.com/thread.jspa?messageID=203510#203510
	if metadata.ETag != nil && fileMD5 == *metadata.ETag {
		return true, nil
	}
	// Case: Server-computed MD5 hash exists
	if objectMD5, ok := metadata.Metadata["fcopy-MD5"]; ok {
		return *objectMD5 == fileMD5, nil
	}

	// Case: No existing hash, or MD5 hash is not the same
	var objectMD5 string
	var err error
	if utils.ObjectIsMultipart(*object.ETag) {
		// TODO: Compute multipart hash for object
	} else {
		objectMD5, err = utils.ObjectMD5Checksum(object)
		if err != nil {
			return false, err
		}
	}
	return objectMD5 == fileMD5, nil
}
