package main

import (
	"bucket-finder/internal/clients"
	"context"
	"fmt"
	"log"
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

	bucketCountEnv, err := strconv.Atoi(os.Getenv("BUCKET_COUNT"))
	if err != nil || bucketCountEnv <= 0 {
		log.Fatal("BUCKET_COUNT must be a positive integer")
	}
	bucketCount := int32(bucketCountEnv)

	address := fmt.Sprintf("%s:%s", os.Getenv("ASSIGNMENT_SERVICE_GRPC_SERVER_ADDRESS"), os.Getenv("ASSIGNMENT_SERVICE_GRPC_SERVER_PORT"))
	client, err := clients.NewGrpcAssignmentClient(address)
	if err != nil {
		log.Fatal(err)
	}

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

	userIds := make([]string, userCount)
	for i := 0; i < userCount; i++ {
		userIds[i] = fmt.Sprintf("user-%d", i)
	}

	variantUserIds := make(map[string][]string)

	ctx := context.Background()
	experimentsAndContext, err := client.GetExperimentsAndVariantsForUsers(ctx, userIds)
	if err != nil {
		log.Fatal(err)
	}

	for userId, experiments := range experimentsAndContext {
		for _, variant := range experiments {
			variantUserIds[variant] = append(variantUserIds[variant], userId)
		}
	}

	for variantKey, userIds := range variantUserIds {
		writeVariantUserIdsFile(variantKey, userIds)
	}

}

func writeVariantUserIdsFile(variantKey string, userIds []string) {
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

func loadEnv() {
	_ = godotenv.Load("../../infrastructure/.env")
}
