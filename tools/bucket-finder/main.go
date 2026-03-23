package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"maps"
	"math/big"
	"os"
	"slices"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	salt := os.Getenv("SALT_VALUE")
	bucketCount := int32(100)

	bucketToUserId := make(map[int32][]int, int(bucketCount))

	for i := 0; i < 100000; i++ {
		bucket := GetBucketFor(strconv.Itoa(i), bucketCount, salt)
		bucketToUserId[bucket] = append(bucketToUserId[bucket], i)
	}

	sortedBucketIds := slices.Sorted(maps.Keys(bucketToUserId))

	for _, bucketId := range sortedBucketIds {
		fmt.Println("Bucket", bucketId, "has user ids:", bucketToUserId[bucketId])
	}
}

func GetBucketFor(userId string, bucketCount int32, salt string) int32 {
	hash := createMD5For(userId, salt)
	hashHex := hex.EncodeToString(hash[:])

	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	bucket := new(big.Int)
	bucket.Mod(hashInt, big.NewInt(int64(bucketCount)))

	return int32(bucket.Int64())
}

func createMD5For(userId, salt string) [16]byte {
	toHash := fmt.Sprintf("%s:%s", salt, userId)
	return md5.Sum([]byte(toHash))
}

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}
