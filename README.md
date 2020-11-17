# File Backup

This program recursively traverses the files of a directory and makes a backup to the cloud.

## Examples

```
> fcopy backup <directory_name> <bucket_name>
> fcopy backup --directory <directory> --bucket <bucket> 
```

```
> fcopy restore <bucket_name> <directory_name>
> fcopy restore --directory <directory> --bucket <bucket> 
```

## Usage

```
go install github.com/hollowsunsets/fcopy
```

```
export AWS_PROFILE=<profile>
export AWS_REGION=us-east-1
```



## Notes


S3 objects will have an etag that is not an MD5 digest if they are uploaded as multipart files.
See https://stackoverflow.com/questions/12186993/what-is-the-algorithm-to-compute-the-amazon-s3-etag-for-a-file-larger-than-5gb for more information.