---
sidebar_position: 4
---

# Data Cooking Service

The data cooking service consumes events from the kafka `events` topic and then joins the event with context for each experiment it is relevant to.

Where an uncooked events is:

```json
{
  "id": "evt_01HXYZABCDE1234567890",
  "event_key": "checkout_completed",
  "user_details": { "id": "user_8f3c2a91" },
  "sent_at": "2026-06-19T14:22:08.512Z",
  "received_at": "2026-06-19T14:22:08.874Z",
  "properties": {
    "cart_total": { "data_type": "float", "value": 49.99 },
    "item_count": { "data_type": "int", "value": 3 },
    "currency": { "data_type": "string", "value": "GBP" }
  }
}
```

and an uncooked `experiment_exposure` event is (Note the addition of the experiment_key):

```json
{
  "id": "evt_01HXYZEXPOSURE0000001",
  "event_key": "experiment_exposure",
  "experiment_key": "checkout_redesign_v2",
  "user_details": { "id": "user_8f3c2a91" },
  "sent_at": "2026-06-19T14:22:08.512Z",
  "received_at": "2026-06-19T14:22:08.874Z",
  "properties": {}
}
```

Where the same events After the cooking process would be written to the database like so (they would be as a database row, but they are shown here as JSON for readability).

```json
{
  "experiment_key": "checkout_redesign_v2",
  "variant_key": "treatment_a",
  "event_key": "checkout_completed",
  "user_id": "user_8f3c2a91",
  "sent_at": "2026-06-19 14:22:08",
  "received_at": "2026-06-19 14:22:08",
  "string_properties": { "currency": "GBP" },
  "int_properties": { "item_count": 3 },
  "float_properties": { "cart_total": 49.99 },
  "is_aa": false
}
```

```json
{
  "experiment_key": "checkout_redesign_v2",
  "variant_key": "treatment_a",
  "event_key": "experiment_exposure",
  "user_id": "user_8f3c2a91",
  "sent_at": "2026-06-19 14:22:08",
  "received_at": "2026-06-19 14:22:08",
  "string_properties": {},
  "int_properties": {},
  "float_properties": {},
  "is_aa": false
}
```

## Process Diagram

The data cooking service's place within the event ingestion pipeline can be seen in the below diagram of the event ingestion pipeline.

![Assignment Diagram](/img/event-ingestion.png)

## Microbatching

The data cooking service uses a custom micro-batching golang library found in `libs/prismmicrobatcher/microbatchingService.go` to read in batches of events from kafka into an internal buffer to then cook events. This approach has been adopted as it matches clickhouse batch write model - the database to which Prism writes cooked events too.

The Kafka library used, KGO, would usually perform commits automatically as soon as they're read from Kafka by the consumer. However, this is risky, as it fails to consider the case where the service fails. To get around this, a manual commit method has been implemented in that the commit to the read head will not take place until the events are safe inside the database.

The size of a micro batch is tunable via environment variables, however, it's not recommended to drop below 10,000 per microbatch.

## Internal Caching

In an effort to reduce load on the service, the data-cooking-service implements free internal caches, which have a lifetime of the micro-batch itself.

| Cache                       | Key                                                                             | Purpose                                                                                                                                                                           |
| --------------------------- | ------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| User ID to Bucket           | User ID                                                                         | Prevents the need to rehash user IDs.                                                                                                                                             |
| Experiment Assignment       | Composite key of the bucket number and the 24-hour truncated event sent at time | Builds a map of the active experiments for a bucket on a certain day, Reducing the number of required GRPC calls.                                                                 |
| Enriched Experiment Details | Experiment key                                                                  | Stores a map of experiment key to experiment details to prevent the need for further gRPC callouts to the Experimentation Service in order to gather enriched experiment context. |
