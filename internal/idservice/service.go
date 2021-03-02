package idservice

import (
	"context"
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
