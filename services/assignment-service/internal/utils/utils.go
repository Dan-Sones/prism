package utils

import (
	"errors"
	"log"
	"os"
	"strconv"
)

func getBucketCount() (int32, error) {
	bCount := os.Getenv("BUCKET_COUNT")
	if bCount == "" {
		return 0, errors.New("BUCKET_COUNT environment variable is not set")
	}

	bCountInt64, err := strconv.ParseInt(bCount, 10, 32)
	if err != nil {
		return 0, errors.New("BUCKET_COUNT environment variable could not be parsed to int32")
	}

	return int32(bCountInt64), nil
}

func GetBucketConfig() (salt string, bucketCount int32) {
	s := os.Getenv("SALT_VALUE")
	bCount, err := getBucketCount()
	if err != nil {
		log.Fatalf(err.Error())
	}
	return s, bCount
}
