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
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"k8s.io/klog/v2"
)

const (
	defaultRegion  = "us-east-1"
	requestTimeout = 15 * time.Second
)

type S3Client struct {
	S3Service *s3.Client
}

func InitS3Client(accessKeyID, secretAccessKey, serviceEndpoint, region string, certData []byte, enableDebug bool) (*S3Client, error) {
	httpClient := http.Client{
		Timeout: requestTimeout,
	}

	skipTLSValidation := false
	if strings.HasPrefix(serviceEndpoint, "https") && len(certData) == 0 {
		skipTLSValidation = true
	}
	if len(certData) > 0 || skipTLSValidation {
		httpClient.Transport = configureTLSTransport(certData, skipTLSValidation)
	}

	if region == "" {
		region = defaultRegion
	}

	// Custom AWS config with static credentials
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, ""))),
		config.WithHTTPClient(&httpClient),
	)
	if err != nil {
		return nil, err
	}

	s3Service := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3Client{
		S3Service: s3Service,
	}, nil
}

func (s *S3Client) CreateBucket(bucketName string) error {
	klog.InfoS("Starting bucket creation", "bucketName", bucketName)

	bucketParams := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	_, err := s.S3Service.CreateBucket(context.TODO(), bucketParams)
	if err != nil {
		// Check for specific S3 errors
		if strings.Contains(err.Error(), "BucketAlreadyExists") || strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			klog.InfoS("Bucket already exists or owned by you", "bucketName", bucketName)
			return nil // Return nil to indicate no error for idempotency
		}
		return err
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
