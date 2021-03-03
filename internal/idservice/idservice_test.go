package idservice

import (
	"context"
	"errors"
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
		cps[i] = types.CommonPrefix{Prefix: aws.String(p)}
	}
	return &s3.ListObjectsV2Output{CommonPrefixes: cps}
}
func TestListObjects(t *testing.T) {

	cases := []struct {
		comment  string
		prefixes []string
		expected []*string
	}{
		{
			comment:  "No prefixes",
			prefixes: []string{},
			expected: []*string{},
		},
		{
			comment:  "One prefix",
			prefixes: []string{"resources/pid1/"},
			expected: []*string{aws.String("pid1")},
		},
		{
			comment:  "Five prefixes",
			prefixes: []string{"resources/pid1/", "resources/pid2/", "resources/pid3/", "resources/pid4/", "resources/pid5/"},
			expected: []*string{aws.String("pid1"), aws.String("pid2"), aws.String("pid3"), aws.String("pid4"), aws.String("pid5")},
		},
	}

	for _, c := range cases {
		mockS3 := new(MyMockedS3)
		input := s3.ListObjectsV2Input{
			Bucket:    aws.String("myBucket"),
			Delimiter: aws.String("/"),
			Prefix:    aws.String("resources/"),
		}
		mockS3.On("ListObjectsV2", context.TODO(), &input).Return(makeOutput(c.prefixes), nil)
		t.Run(c.comment,
			func(t *testing.T) {

				uut, err := NewIDService(S3Client(mockS3))
				if err != nil {
					t.Errorf("error: %w", err)
				}

				actual, _ := uut.List()

				mockS3.AssertExpectations(t)
				assert.ElementsMatch(t, c.expected, actual)
			})
	}
}

func TestGetResourceId(t *testing.T) {
	cases := []struct {
		comment  string
		key      string
		expected int
		err      error
	}{
		// {
		// 	comment:  "One value",
		// 	key:      "resources/pid1/1",
		// 	expected: 1,
		// },
		{
			comment:  "No value",
			key:      "",
			expected: 1,
			err:      errors.New("Could not obtain resource ID"),
		},
	}
	for _, c := range cases {
		mockS3 := new(MyMockedS3)
		input := s3.ListObjectsV2Input{
			Bucket:    aws.String("myBucket"),
			Delimiter: aws.String("/"),
			Prefix:    aws.String("resources/pid1/"),
		}
		mockS3.On("ListObjectsV2", context.TODO(), &input).Return(makeOutput([]string{c.key}), nil)
		t.Run(c.comment,
			func(t *testing.T) {

				uut, err := NewIDService(S3Client(mockS3))
				if err != nil {
					t.Errorf("error: %w", err)
				}

				actual, err := uut.GetResourceID("pid1")
				mockS3.AssertExpectations(t)
				if c.err == nil {
					if err != nil {
						t.Errorf("error: %w", err)
					}
					assert.Equal(t, c.expected, actual)
				} else {
					assert.Equal(t, c.err, err)
				}

			})
	}
}
