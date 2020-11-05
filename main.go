package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an action for the tool to take")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "backup":
		fmt.Println("backup")
	case "restore":
		fmt.Println("restore")
	default:
		fmt.Printf("Unsupported action: %s\n", os.Args[1])
		os.Exit(1)
	}

	bucketName := flag.String("bucketName", "", "Name of S3 bucket")
	directoryName := flag.String("directoryName", "", "Name of file directory")
	flag.Parse()

	if *bucketName == "" || *directoryName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func backup(directoryName string, bucketName string) {

}

func restore(bucketName string, directoryName string) {

}
