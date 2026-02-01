package clients

import (
	"fmt"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

func GetKafkaClient() (*kgo.Client, error) {
	broker := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_BOOTSTRAP_SERVER_HOST"), os.Getenv("KAFKA_BOOTSTRAP_SERVER_PORT"))

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
		kgo.ConsumerGroup(os.Getenv("ASSIGNMENT_SERVICE_KAFKA_CONSUMER_GROUP_ID")),
		kgo.ConsumeTopics(os.Getenv("KAFKA_CACHE_INVALIDATIONS_TOPIC")),
	)
	if err != nil {
		return nil, err
	}

	return cl, nil
}
