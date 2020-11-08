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
```