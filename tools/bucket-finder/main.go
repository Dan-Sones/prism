package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"golang.design/x/clipboard"
)

func main() {
	loadEnv()

	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	salt := os.Getenv("SALT_VALUE")
	bucketCount := int32(100)

	assignementServiceGetAssignmentsUrl := fmt.Sprintf("http://%s:%s/api/assignments/", "localhost", os.Getenv("ASSIGNMENT_SERVICE_HTTP_PORT"))

	bucketToUserId := make(map[int32][]int, int(bucketCount))

	fmt.Println("Please enter the number of users to bucket (e.g. 10000000):")
	var input string
	fmt.Scanln(&input)

	userCount, err := strconv.Atoi(input)
	if err != nil || userCount <= 0 {
		fmt.Println("Invalid number of users. Please enter a positive integer.")
		return
	}

	fmt.Println("Please enter the buckets in the format 1,2,3 (e.g. 0,1,2):")
	var bucketInput string
	fmt.Scanln(&bucketInput)

	buckets := strings.Split(bucketInput, ",")
	for _, bucketStr := range buckets {
		bucketId, err := strconv.Atoi(bucketStr)
		if err != nil || bucketId < 0 || bucketId >= int(bucketCount) {
			fmt.Printf("Invalid bucket id '%s'. Please enter numbers between 0 and 99.\n", bucketStr)
			return
		}
	}

	for i := 0; i < userCount; i++ {
		bucket := GetBucketFor(strconv.Itoa(i), bucketCount, salt)
		bucketToUserId[bucket] = append(bucketToUserId[bucket], i)
	}

	var userIds []int
	for _, bucketStr := range buckets {
		bucketId, _ := strconv.Atoi(bucketStr)
		userIds = append(userIds, bucketToUserId[int32(bucketId)]...)
	}

	variantUserIds := make(map[string][]int)

	for _, userId := range userIds {
		address := fmt.Sprintf("%s%d", assignementServiceGetAssignmentsUrl, userId)
		fmt.Println(address)
		res, err := http.Get(address)
		if err != nil {
			fmt.Printf("Failed to get assignments for user id %d: %v\n", userId, err)
			return
		}

		if res.StatusCode == http.StatusOK {
			var assignments map[string]string
			err = json.NewDecoder(res.Body).Decode(&assignments)
			if err != nil {
				fmt.Printf("Failed to decode response for user id %d: %v\n", userId, err)
				return
			}

			for _, variant := range assignments {
				variantUserIds[variant] = append(variantUserIds[variant], userId)
			}

		} else {
			fmt.Printf("Received non-OK response for user id %d: %s\n", userId, res.Status)
		}
	}

	for variantKey, userIds := range variantUserIds {
		writeVariantUserIdsFile(variantKey, userIds)
	}

}

func writeVariantUserIdsFile(variantKey string, userIds []int) {
	data, err := yaml.Marshal(userIds)
	if err != nil {
		fmt.Printf("Error marshalling user ids for variant %s: %v\n", variantKey, err)
		return
	}

	fileName := fmt.Sprintf("../experiment-simulator/resources/variant-%s-users.yml", variantKey)
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file for variant %s: %v\n", variantKey, err)
		return
	}

	fmt.Printf("Written %d user ids for variant %s to %s\n", len(userIds), variantKey, fileName)
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
