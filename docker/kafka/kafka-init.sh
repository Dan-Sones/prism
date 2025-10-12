#!/bin/bash

create_topic() {
  local topic=$1
  /usr/bin/kafka-topics --bootstrap-server localhost:9092 \
    --create --if-not-exists \
    --topic "$topic" \
    --replication-factor 1 \
    --partitions 1
}

until /usr/bin/kafka-topics --bootstrap-server localhost:9092 --list; do
  echo "Waiting for Kafka to be ready..."
  sleep 5
done

create_topic new-topic