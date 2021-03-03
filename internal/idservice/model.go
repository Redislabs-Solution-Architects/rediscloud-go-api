package idservice

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//S3ClientType represents the capabilites of the s3 client
type S3ClientType s3.ListObjectsV2APIClient

//IDService is the implementation of the IDService
type IDService struct {
	S3Client  S3ClientType
	bucket    string
	delimiter string
	prefix    string
}

//NewIDService returns an IDService configured with default parameters
func NewIDService(opts ...ServiceOption) (svc *IDService, err error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return
	}

	service := &IDService{
		S3Client:  s3.NewFromConfig(cfg),
		bucket:    "myBucket",
		delimiter: "/",
		prefix:    "resources",
	}
	for _, opt := range opts {
		opt(service)
	}

	return service, err
}

//ServiceOption is a function type that will modify the provided idservice
type ServiceOption func(*IDService)

//S3Client will set the idservice's S3Client to the provided client
func S3Client(client S3ClientType) ServiceOption {
	return func(i *IDService) {
		i.S3Client = client
	}
}

//Bucket will modify the bucket
func Bucket(bucket string) ServiceOption {
	return func(i *IDService) {
		i.bucket = bucket
	}
}

//Delimiter will modify the bucket
func Delimiter(delimiter string) ServiceOption {
	return func(i *IDService) {
		i.delimiter = delimiter
	}
}

//Prefix will modify the bucket
func Prefix(prefix string) ServiceOption {
	return func(i *IDService) {
		i.prefix = prefix
	}
}

//Interface determines the methods available
type Interface interface {
	// Put(primarId string, resourceId int) error
	GetResourceID(primaryID string) (int, error)
	List() ([]string, error)
	// Delete (primaryId string) error
}
