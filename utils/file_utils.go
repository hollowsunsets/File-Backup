package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"math"
	"os"
	"strings"
)

const chunkSize = 8 * 1024 * 1024
func FileMD5Checksum(file* os.File) (string, error) {
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := info.Size()

	blocks := uint64(math.Ceil(float64(fileSize) / float64(chunkSize)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		chunksSoFar := int64(i * chunkSize)
		blockSize := int(math.Min(chunkSize, float64(fileSize - chunksSoFar)))
		buffer := make([]byte, blockSize)

		file.Read(buffer)
		io.WriteString(hash, string(buffer))
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetObject(bucketName string, fileName string) (*s3.GetObjectOutput, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
	}

	result, err := svc.GetObject(params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetObjectMetadata(bucketName string, fileName string) (*s3.HeadObjectOutput, error) {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	output, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key: aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}
	return output, nil
}

func ObjectIsMultipart(etag string) bool {
	return len([]rune(etag)) != 32 && strings.HasSuffix(etag,"-#")
}

func BucketExists(bucketName string) bool {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	}
	_, err := svc.HeadBucket(input)
	if err != nil {
		return false
	}
	return true
}

func CreateBucket(bucketName string) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	_, err := svc.CreateBucket(input)
	if err != nil {
		return err
	}
	return nil
}

func ObjectMD5Checksum(object *s3.GetObjectOutput) (string, error) {
	return "", nil
}