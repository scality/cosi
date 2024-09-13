package driver

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	s3client "github.com/scality/cosi/pkg/util/s3client"
	mock_s3iface "github.com/scality/cosi/pkg/util/s3client/mock" // Import the generated mock
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/kubernetes"
	cosispec "sigs.k8s.io/container-object-storage-interface-spec"
)

// Test successful bucket creation
func TestDriverCreateBucket_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock S3 client
	mockS3 := mock_s3iface.NewMockS3API(ctrl)

	// Mock the S3 CreateBucket call
	mockS3.EXPECT().CreateBucket(gomock.Any()).Return(nil, nil).Times(1)

	// Mock the initializeObjectStorageProviderClients function to return the mock S3 client
	initializeS3Client = func(ctx context.Context, clientset *kubernetes.Clientset, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	// Call DriverCreateBucket
	resp, err := server.DriverCreateBucket(context.TODO(), req)

	// Assert that no error occurred and the response is as expected
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-bucket", resp.BucketId)
}

// Test S3 client initialization failure
func TestDriverCreateBucket_InitClientFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock initializeObjectStorageProviderClients to return an error
	initializeS3Client = func(ctx context.Context, clientset *kubernetes.Clientset, parameters map[string]string) (*s3client.S3Client, error) {
		return nil, errors.New("failed to initialize S3 client")
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	// Call DriverCreateBucket
	resp, err := server.DriverCreateBucket(context.TODO(), req)

	// Assert the error is codes.Internal and no response
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

// Test bucket already exists case
func TestDriverCreateBucket_BucketAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock S3 client
	mockS3 := mock_s3iface.NewMockS3API(ctrl)

	// Mock the S3 CreateBucket call to return bucket already exists error
	awsErr := awserr.New(s3.ErrCodeBucketAlreadyExists, "Bucket already exists", nil)
	mockS3.EXPECT().CreateBucket(gomock.Any()).Return(nil, awsErr).Times(1)

	// Mock the initializeObjectStorageProviderClients function to return the mock S3 client
	initializeS3Client = func(ctx context.Context, clientset *kubernetes.Clientset, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	// Call DriverCreateBucket
	resp, err := server.DriverCreateBucket(context.TODO(), req)

	// Assert the error is codes.AlreadyExists
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

// Test bucket already owned by you case
func TestDriverCreateBucket_BucketAlreadyOwnedByYou(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock S3 client
	mockS3 := mock_s3iface.NewMockS3API(ctrl)

	// Mock the S3 CreateBucket call to return bucket already owned by you error
	awsErr := awserr.New(s3.ErrCodeBucketAlreadyOwnedByYou, "Bucket already owned by you", nil)
	mockS3.EXPECT().CreateBucket(gomock.Any()).Return(nil, awsErr).Times(1)

	// Mock the initializeObjectStorageProviderClients function to return the mock S3 client
	initializeS3Client = func(ctx context.Context, clientset *kubernetes.Clientset, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	// Call DriverCreateBucket
	resp, err := server.DriverCreateBucket(context.TODO(), req)

	// Assert the error is codes.AlreadyExists
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}
