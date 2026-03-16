package clients

import (
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func GetKafkaClient() (*kgo.Client, error) {
	broker := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_BOOTSTRAP_SERVER_HOST"), os.Getenv("KAFKA_BOOTSTRAP_SERVER_PORT"))

	cl, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
		kgo.ConsumerGroup(os.Getenv("DATA_COOKING_SERVICE_KAFKA_CONSUMER_GROUP_ID")),
		kgo.ConsumeTopics(os.Getenv("KAFKA_EVENTS_TOPIC")),
		kgo.FetchMinBytes(1000000),
		kgo.FetchMaxWait(2*time.Second),
		kgo.SessionTimeout(6*time.Second),
		kgo.HeartbeatInterval(2*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return cl, nil
}
