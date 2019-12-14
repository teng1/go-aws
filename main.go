package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func listBuckets() {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")},
	)

	// Create S3 service client
	svc := s3.New(sess)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	// return result.Buckets

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}

}

func listObjects(bucket string) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-2")})
	if err != nil {
		exitErrorf("Unable to make session", err)
	}
	// // Create S3 service client
	svc := s3.New(sess)

	params := &s3.ListObjectsV2Input{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(100),
	}
	truncatedListing := true

	for truncatedListing {
		resp, err := svc.ListObjectsV2(params)
		if err != nil {
			// Print the error.
			exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
			return
		}
		// Get all files
		for _, item := range resp.Contents {
			fmt.Println(*item.Key)
		}
		// Set continuation token
		params.ContinuationToken = resp.NextContinuationToken
		truncatedListing = *resp.IsTruncated
	}
}

func main() {
	listBuckets()
	listObjects("my-bucket")

}
