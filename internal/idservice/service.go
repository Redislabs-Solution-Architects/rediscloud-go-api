package idservice

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//List returns a slice of the pids known, or an error indicating a problem
func (svc IDService) List() (pids []*string, err error) {
	in := &s3.ListObjectsV2Input{
		Bucket:    aws.String(svc.bucket),
		Delimiter: aws.String(svc.delimiter),
		Prefix:    aws.String(svc.prefix + svc.delimiter),
	}
	output, err := svc.S3Client.ListObjectsV2(context.TODO(), in)
	if err != nil {
		return
	}

	for _, prefix := range output.CommonPrefixes {
		p := aws.ToString(prefix.Prefix)
		pid := strings.Split(p, svc.delimiter)[1]
		pids = append(pids, &pid)
	}

	return
}

//GetResourceID returns the resource id for the given primary id, or an error if that's not possible
func (svc IDService) GetResourceID(primaryID string) (rid int, err error) {
	in := &s3.ListObjectsV2Input{
		Bucket:    aws.String(svc.bucket),
		Delimiter: aws.String(svc.delimiter),
		Prefix:    aws.String(svc.prefix + svc.delimiter + primaryID + svc.delimiter),
	}
	output, err := svc.S3Client.ListObjectsV2(context.TODO(), in)
	if err != nil {
		return
	}

	defer func() {
		if v := recover(); v != nil {
			err = errors.New("Could not obtain resource ID")
		}
	}()
	r := strings.Split(*output.CommonPrefixes[0].Prefix, svc.delimiter)[2]
	rid, err = strconv.Atoi(r)
	return
}
