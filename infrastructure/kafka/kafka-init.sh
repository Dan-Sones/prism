#!/bin/bash

create_topic() {
  local topic=$1
  local partitions=${2:-1}
  local bootstrap=${KAFKA_BOOTSTRAP_SERVER:-localhost:9092}
  /usr/bin/kafka-topics --bootstrap-server "$bootstrap" \
    --create --if-not-exists \
    --topic "$topic" \
    --replication-factor 1 \
    --partitions "$partitions"
}
# TODO: the replication-factor and partitions need to be more than this for a production use most likely (?)
# partitions default at one, need to spin up more brokers for increased replication factor 

KAFKA_BOOTSTRAP_SERVER=${KAFKA_BOOTSTRAP_SERVER:-localhost:9092}

until /usr/bin/kafka-topics --bootstrap-server "$KAFKA_BOOTSTRAP_SERVER" --list; do
  echo "Waiting for Kafka to be ready at $KAFKA_BOOTSTRAP_SERVER..."
  sleep 5
done

create_topic "${KAFKA_CACHE_INVALIDATIONS_TOPIC:-assignment-cache-invalidations}"
create_topic "${KAFKA_EVENTS_TOPIC:-events}"
