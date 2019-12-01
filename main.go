package main

import (
	"errors"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var s3svc = s3.New(session.New(&aws.Config{
	Region: aws.String("eu-west-1")}),
)

func findStateBucket() (stateBucket string, err error) {
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Error("Failed to retrieve buckets: ", err)
	}
	for _, bucket := range result.Buckets {
		stateBucketTrue := strings.Contains(aws.StringValue(bucket.Name), "state.terragrunt")
		if stateBucketTrue {
			log.Debug("Found state bucket: ", aws.StringValue(bucket.Name))
			return aws.StringValue(bucket.Name), nil
		} else {
			continue
		}
	}
	return "", errors.New("unable to find state bucket")
}

func getVersioningStatus(bucket string) (bucketVersioningStatus string, err error) {
	input := &s3.GetBucketVersioningInput{
		Bucket: aws.String(bucket),
	}

	result, err := s3svc.GetBucketVersioning(input)
	if err != nil {
		return "", error(err)
	}
	return result.String(), nil
}

func getVersioningEnabled(status string) (enabled bool) {
	if strings.Contains(strings.ToLower(status), "enabled") {
		return true
	} else {
		return false
	}
}

func main() {
	appLogLevel, set := os.LookupEnv("LOG_LEVEL")
	if !set {
		log.SetLevel(log.InfoLevel)
	} else {
		logLevel, err := log.ParseLevel(appLogLevel)
		if err != nil {
			log.Error(err)
		}
		log.SetLevel(logLevel)
	}

	stateBucket, err := findStateBucket()
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("State bucket found: ", stateBucket)

	versioningStatus, err := getVersioningStatus(stateBucket)
	if err != nil {
		log.Error(err)
	} else {
		log.Debug(stateBucket, versioningStatus)
	}

	log.Infof("state bucket: %s\t versioned: %v", stateBucket, getVersioningEnabled(versioningStatus))
}
