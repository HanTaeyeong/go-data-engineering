package functions

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(fileName string, data string) {

	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRETE_KEY")

	var options = s3.Options{
		Region:      "ap-northeast-2",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
	}

	var client = s3.New(options, func(o *s3.Options) {
		o.Region = "ap-northeast-2"
	})

	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("farmers-realtime-data"),
		Key:    aws.String(fileName),
		Body:   strings.NewReader(data),
	})

	if err != nil {
		panic(err)
	}

	//fmt.Print("uploadSuccess", response)
}
