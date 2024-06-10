package internal

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var Uploader *s3manager.Uploader

func NewAWS() {

	var region string = ""    // AWS Region
	var accessKey string = "" // access key
	var secretKey string = "" // secret key

	awsSession, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				Credentials: credentials.NewStaticCredentials(
					accessKey,
					secretKey,
					"",
				),
			},
		})

	if err != nil {
		panic(err)
	}

	Uploader = s3manager.NewUploader(awsSession)
}

type Result struct {
	Value string
	Err   error
}

func UploadImage(file *multipart.FileHeader) <-chan Result {
	ch := make(chan Result)

	go func() {
		defer close(ch)
		src, err := file.Open()
		if err != nil {
			return
		}
		defer src.Close()

		var bucketName string = "rohan-blog-bucket" // bucket name

		_, err = Uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(file.Filename),
			Body:   src, // add file body here
		})
		if err != nil {
			ch <- Result{Value: "", Err: err}
			return
		}

		url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, file.Filename)

		ch <- Result{Value: url, Err: nil}
	}()

	return ch
}
