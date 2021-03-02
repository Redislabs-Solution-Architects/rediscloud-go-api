package idservice

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MyMockedS3 struct {
	mock.Mock
}

func (m *MyMockedS3) ListObjectsV2(ctx context.Context, in *s3.ListObjectsV2Input, optFuns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	// convert to variadic type
	var args mock.Arguments
	if len(optFuns) > 0 {
		var ofs []interface{} = make([]interface{}, len(optFuns))
		for i, of := range optFuns {
			ofs[i] = of
		}
		args = m.Called(ctx, in, ofs)
	} else {
		args = m.Called(ctx, in)
	}
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func makeOutput(prefixes []string) *s3.ListObjectsV2Output {
	cps := make([]types.CommonPrefix, len(prefixes))
	for i, p := range prefixes {
		cps[i] = types.CommonPrefix{&p}
	}
	return &s3.ListObjectsV2Output{CommonPrefixes: cps}
}
func TestListObjects(t *testing.T) {
	mockS3 := new(MyMockedS3)
	input := s3.ListObjectsV2Input{
		Bucket:    aws.String("myBucket"),
		Delimiter: aws.String("/"),
		Prefix:    aws.String("resources/"),
	}
	output := s3.ListObjectsV2Output{
		CommonPrefixes: []types.CommonPrefix{{Prefix: aws.String("resources/pid1/")}},
	}
	//prefixes := []string{"resources/pid1/"}

	mockS3.On("ListObjectsV2", context.TODO(), &input).Return(&output, nil)

	uut, err := NewIDService(S3Client(mockS3))
	if err != nil {
		t.Errorf("error: %w", err)
	}

	actual, _ := uut.List()
	expected := []*string{aws.String("pid1")}
	mockS3.AssertExpectations(t)
	assert.ElementsMatch(t, expected, actual)
}
