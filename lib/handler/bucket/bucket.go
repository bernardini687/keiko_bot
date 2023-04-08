package bucket

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ID struct {
	Bucket *string
	Key    *string
}

type NotFoundError struct {
	Key string
}

func (err *NotFoundError) Error() string {
	return fmt.Sprintf("no data at `%s`", err.Key)
}

func NewID(bucket, key string) *ID {
	return &ID{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
}

func NewSession() *session.Session {
	return session.Must(session.NewSession())
}

func NewClient(sess *session.Session) *s3.S3 {
	return s3.New(sess)
}

func NewUploader(sess *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(sess)
}

func LookupKey(client *s3.S3, id *ID) (bool, error) {
	// HeadObject would be better but it can't report the exact reason of the error and
	// we want to be sure it is a "no such key" one
	_, err := client.GetObject(&s3.GetObjectInput{
		Bucket: id.Bucket,
		Key:    id.Key,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			log.Printf("4/6) aerr code: %#v\n", aerr.Code())
			log.Printf("5/7) err code no such key: %#v\n", s3.ErrCodeNoSuchKey)

			if aerr.Code() == s3.ErrCodeNoSuchKey {
				return false, nil
			}
			return false, aerr
		} else {
			return false, err
		}
	}

	return true, nil
}

func GetContent(client *s3.S3, id *ID) (string, error) {
	res, err := client.GetObject(&s3.GetObjectInput{
		Bucket: id.Bucket,
		Key:    id.Key,
		// TODO: what can i gain from specifying the content type?
	})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func GetContentFromKey(key string) (string, error) {
	id := NewID(os.Getenv("BUCKET_NAME"), key)

	sess := NewSession()
	client := NewClient(sess)
	found, err := LookupKey(client, id)
	if err != nil {
		return "", err
	}
	if !found {
		return "", &NotFoundError{key}
	}

	res, err := client.GetObject(&s3.GetObjectInput{
		Bucket: id.Bucket,
		Key:    id.Key,
		// TODO: what can i gain from specifying the content type?
	})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func UploadContents(uploader *s3manager.Uploader, id *ID, contents string) error {
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: id.Bucket,
		Key:    id.Key,
		Body:   strings.NewReader(contents),
	})
	if err != nil {
		return err
	}

	return nil
}
