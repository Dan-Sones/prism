package prismhash

import (
	"log"
	"os"
	"strconv"
)

func GetBucketConfig() (salt string, bucketCount int32) {
	salt = os.Getenv("SALT_VALUE")
	bCountStr := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.ParseInt(bCountStr, 10, 32)
	if err != nil {
		log.Fatal("BUCKET_COUNT environment variable could not be parsed to int32")
	}
	return salt, int32(bCount)
}
