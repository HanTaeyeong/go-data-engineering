package functions

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(fileName string, data string) {

	accessKey := os.Getenv("accessKey")
	secretKey := os.Getenv("secretKey")

	var options = s3.Options{
		Region:      "ap-northeast-2",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	var client = s3.New(options, func(o *s3.Options) {
		o.Region = "ap-northeast-2"
	})

	response, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("farmers-realtime-data"),
		Key:    aws.String(fileName),
		Body:   strings.NewReader(data),
	})

	if err != nil {
		panic(err)
	}

	fmt.Print("uploadSuccess", response)
}
