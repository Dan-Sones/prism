package clients

import (
	"fmt"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

func NewKafkaProducer() (*kgo.Client, error) {
	broker := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_BOOTSTRAP_SERVER_HOST"), os.Getenv("KAFKA_BOOTSTRAP_SERVER_PORT"))
	client, err := kgo.NewClient(
		kgo.SeedBrokers(broker),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
