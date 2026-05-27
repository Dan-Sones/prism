package clients

import (
	"fmt"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

func GetKafkaProducerClient() (*kgo.Client, error) {
	broker := fmt.Sprintf("%s:%s",
		os.Getenv("KAFKA_BOOTSTRAP_SERVER_HOST"),
		os.Getenv("KAFKA_BOOTSTRAP_SERVER_PORT"),
	)

	return kgo.NewClient(
		kgo.SeedBrokers(broker),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordDeliveryTimeout(30*time.Second),
	)
}
