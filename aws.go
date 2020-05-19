package main

import (
	"bytes"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
	"github.com/joho/godotenv"
)

var (
	AWS_S3_REGION         string
	AWS_S3_BUCKET         string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	PORT                  string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AWS_S3_REGION = os.Getenv("AWS_S3_REGION")
	AWS_S3_BUCKET = os.Getenv("AWS_S3_BUCKET")
	AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
	PORT = os.Getenv("PORT")
}

// Amazon struc
type Amazon struct {
	Region    string
	Bucket    string
	AccessID  string
	AccessKey string
}

// Connect func
func (a *Amazon) Connect() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(a.Region),
		Credentials: credentials.NewStaticCredentials(
			a.AccessID,
			a.AccessKey,
			"",
		),
	})

	if err != nil {
		return nil, err
	}

	return session, nil
}

// UploadFileS3 func
func (a *Amazon) UploadFileS3(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	s, err := a.Connect()
	if err != nil {
		return "", err
	}

	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := "picture-" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(a.Bucket),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})

	if err != nil {
		return "", err
	}

	return tempFileName, err
}
