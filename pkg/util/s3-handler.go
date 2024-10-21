package s3client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"k8s.io/klog/v2"
)

const (
	defaultRegion  = "us-east-1"
	requestTimeout = 15 * time.Second
)

// S3Client wraps the AWS SDK's S3 client to allow for custom methods
type S3Client struct {
	S3Service *s3.Client
}

// InitS3Client initializes and returns a new S3Client instance using AWS SDK v2
func InitS3Client(accessKeyID, secretAccessKey, serviceEndpoint, region string, certData []byte, enableDebug bool) (*S3Client, error) {
	loggingLevel := aws.LogOff
	if enableDebug {
		loggingLevel = aws.LogDebug
	}

	httpClient := http.Client{
		Timeout: requestTimeout,
	}

	enableTLS := false
	skipTLSValidation := false
	if strings.HasPrefix(serviceEndpoint, "https") && len(certData) == 0 {
		skipTLSValidation = true
	}
	if len(certData) > 0 || skipTLSValidation {
		enableTLS = true
		httpClient.Transport = configureTLSTransport(certData, skipTLSValidation)
	}

	if region == "" {
		region = defaultRegion
	}

	// Load the AWS configuration
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithHTTPClient(&httpClient),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client using the custom resolver
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(serviceEndpoint)
	})

	return &S3Client{
		S3Service: s3Client,
	}, nil
}

// CreateBucket attempts to create an S3 bucket with the provided name using the S3 client
func (s *S3Client) CreateBucket(bucketName string) error {
	klog.InfoS("Starting bucket creation", "bucketName", bucketName)

	bucketParams := &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(defaultRegion),
		},
	}

	_, err := s.S3Service.CreateBucket(context.TODO(), bucketParams)
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	klog.InfoS("Bucket creation succeeded", "bucketName", bucketName)
	return nil
}

// configureTLSTransport sets up the HTTP transport with TLS support
func configureTLSTransport(certData []byte, skipTLS bool) *http.Transport {
	tlsSettings := &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: skipTLS}

	if len(certData) > 0 {
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certData); !ok {
			klog.Warning("Failed to append provided cert data to the certificate pool")
		}
		tlsSettings.RootCAs = caCertPool
	}

	return &http.Transport{
		TLSClientConfig: tlsSettings,
	}
}
