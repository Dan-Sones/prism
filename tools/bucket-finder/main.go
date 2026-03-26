package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"golang.design/x/clipboard"
	"gopkg.in/yaml.v3"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	salt := os.Getenv("SALT_VALUE")
	bucketCount := int32(100)

	bucketToUserId := make(map[int32][]int, int(bucketCount))

	for i := 0; i < 100000; i++ {
		bucket := GetBucketFor(strconv.Itoa(i), bucketCount, salt)
		bucketToUserId[bucket] = append(bucketToUserId[bucket], i)
	}

	for {
		fmt.Println("Please enter a bucket id between 0 and 99 to see the user ids in that bucket, or type 'exit' to quit:")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		bucketId, err := strconv.Atoi(input)
		if err != nil || bucketId < 0 || bucketId >= int(bucketCount) {
			fmt.Println("Invalid bucket id. Please enter a number between 0 and 99.")
			continue
		}

		userIds := bucketToUserId[int32(bucketId)]
		fmt.Printf("Bucket %d contains user ids: %v\n", bucketId, userIds)
		fmt.Println("Do you want to copy the yml array for these user ids to your clipboard? (y/n)")
		var copyInput string
		fmt.Scanln(&copyInput)

		if copyInput == "y" {

			data, err := yaml.Marshal(userIds)
			if err != nil {
				fmt.Printf("Error marshalling user ids to YAML: %v\n", err)
				continue
			}

			clipboard.Write(clipboard.FmtText, data)
			fmt.Printf("YAML array for bucket %d copied to clipboard.\n", bucketId)
		}
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
