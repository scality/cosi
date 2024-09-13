package driver

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	s3client "github.com/scality/cosi/pkg/util/s3client"
	mock_s3iface "github.com/scality/cosi/pkg/util/s3client/mock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/api/core/v1"                       // For Kubernetes Secret and other core types
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // For metadata such as ObjectMeta
	"k8s.io/apimachinery/pkg/runtime"             // For runtime.Object
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"  // For creating the fake Kubernetes client
	ktesting "k8s.io/client-go/testing" // Aliased to avoid conflict with testing package
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
	initializeS3Client = func(ctx context.Context, clientset kubernetes.Interface, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	// Call DriverCreateBucket
	resp, err := server.DriverCreateBucket(context.TODO(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-bucket", resp.BucketId)
}

func TestDriverCreateBucket_InitClientFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	initializeS3Client = func(ctx context.Context, clientset kubernetes.Interface, parameters map[string]string) (*s3client.S3Client, error) {
		return nil, errors.New("failed to initialize S3 client")
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	resp, err := server.DriverCreateBucket(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Internal, status.Code(err))
}

func TestDriverCreateBucket_BucketAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_s3iface.NewMockS3API(ctrl)

	awsErr := awserr.New(s3.ErrCodeBucketAlreadyExists, "Bucket already exists", nil)
	mockS3.EXPECT().CreateBucket(gomock.Any()).Return(nil, awsErr).Times(1)

	initializeS3Client = func(ctx context.Context, clientset kubernetes.Interface, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	resp, err := server.DriverCreateBucket(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestDriverCreateBucket_BucketAlreadyOwnedByYou(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_s3iface.NewMockS3API(ctrl)

	awsErr := awserr.New(s3.ErrCodeBucketAlreadyOwnedByYou, "Bucket already owned by you", nil)
	mockS3.EXPECT().CreateBucket(gomock.Any()).Return(nil, awsErr).Times(1)

	// Mock the initializeObjectStorageProviderClients function to return the mock S3 client
	initializeS3Client = func(ctx context.Context, clientset kubernetes.Interface, parameters map[string]string) (*s3client.S3Client, error) {
		return &s3client.S3Client{S3Service: mockS3}, nil
	}
	defer func() { initializeS3Client = initializeObjectStorageProviderClients }() // Reset after the test

	server := &provisionerServer{}

	req := &cosispec.DriverCreateBucketRequest{
		Name: "test-bucket",
	}

	resp, err := server.DriverCreateBucket(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.AlreadyExists, status.Code(err))
}

func TestDriverDeleteBucket_Unimplemented(t *testing.T) {
	server := &provisionerServer{}

	req := &cosispec.DriverDeleteBucketRequest{
		BucketId: "test-bucket",
	}

	resp, err := server.DriverDeleteBucket(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Unimplemented, status.Code(err))
}

func TestDriverGrantBucketAccess_Unimplemented(t *testing.T) {
	server := &provisionerServer{}

	req := &cosispec.DriverGrantBucketAccessRequest{
		BucketId: "test-bucket",
	}

	resp, err := server.DriverGrantBucketAccess(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Unimplemented, status.Code(err))
}

func TestDriverRevokeBucketAccess_Unimplemented(t *testing.T) {
	server := &provisionerServer{}

	req := &cosispec.DriverRevokeBucketAccessRequest{
		BucketId: "test-bucket",
	}

	resp, err := server.DriverRevokeBucketAccess(context.TODO(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.Unimplemented, status.Code(err))
}

// Test initializeObjectStorageProviderClients with successful secret fetch
func TestInitializeObjectStorageProviderClients_Success(t *testing.T) {
	// Use a fake Kubernetes client instead of NewForConfigOrDie(nil)
	fakeClientset := fake.NewSimpleClientset()

	// Add a fake secret to the client
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-secret",
			Namespace: "test-namespace",
		},
		Data: map[string][]byte{
			"COSI_S3_ACCESS_KEY_ID":     []byte("access-key"),
			"COSI_S3_ACCESS_SECRET_KEY": []byte("secret-key"),
			"COSI_S3_ENDPOINT":          []byte("http://localhost"),
			"COSI_S3_REGION":            []byte("us-west-1"),
		},
	}

	// Mock the secret retrieval
	fakeClientset.PrependReactor("get", "secrets", func(action ktesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, secret, nil
	})

	// Call initializeObjectStorageProviderClients with the fake client
	parameters := map[string]string{
		"COSI_OBJECT_STORAGE_PROVIDER_SECRET_NAME":      "test-secret",
		"COSI_OBJECT_STORAGE_PROVIDER_SECRET_NAMESPACE": "test-namespace",
	}
	client, err := initializeObjectStorageProviderClients(context.TODO(), fakeClientset, parameters)

	assert.NoError(t, err)
	assert.NotNil(t, client)
}
