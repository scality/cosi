/*
Copyright 2024 Scality, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package s3client

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"k8s.io/klog/v2"
)

const (
	defaultRegion  = "us-east-1"
	requestTimeout = 15 * time.Second
)

// S3Client wraps the s3iface.S3API structure to allow for custom methods
type S3Client struct {
	S3Service s3iface.S3API
}

// InitS3Client initializes and returns a new S3Client instance
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

	awsSession, sessionErr := session.NewSession(
		aws.NewConfig().
			WithRegion(region).
			WithCredentials(credentials.NewStaticCredentials(accessKeyID, secretAccessKey, "")).
			WithEndpoint(serviceEndpoint).
			WithS3ForcePathStyle(true).
			WithMaxRetries(5).
			WithDisableSSL(!enableTLS).
			WithHTTPClient(&httpClient).
			WithLogLevel(loggingLevel),
	)
	if sessionErr != nil {
		return nil, sessionErr
	}

	s3Service := s3.New(awsSession)
	return &S3Client{
		S3Service: s3Service,
	}, nil
}

// CreateBucket attempts to create an S3 bucket with the provided name
func (s *S3Client) CreateBucket(bucketName string) error {
	klog.InfoS("Starting bucket creation", "bucketName", bucketName)

	bucketParams := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	_, creationErr := s.S3Service.CreateBucket(bucketParams)
	if creationErr != nil {
		return creationErr
	}

	klog.InfoS("Bucket creation succeeded", "bucketName", bucketName)
	return nil
}

// configureTLSTransport sets up the HTTP transport with TLS support
func configureTLSTransport(certData []byte, skipTLS bool) *http.Transport {
	tlsSettings := &tls.Config{MinVersion: tls.VersionTLS12, InsecureSkipVerify: skipTLS}

	if len(certData) > 0 {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(certData)
		tlsSettings.RootCAs = caCertPool
	}

	return &http.Transport{
		TLSClientConfig: tlsSettings,
	}
}
