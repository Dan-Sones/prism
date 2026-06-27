---
sidebar_position: 4
---

# Events Service

The Event Service is perhaps the simplest service in Prism's architecture. It has a single responsibility in validating and ingesting events.

## Flow

Events arrive to the events service via a HTTP POST request that looks like the following:

```json
{
  "event_key": "checkout_completed",
  "user_details": { "id": "user_8f3c2a91" },
  "sent_at": "2026-06-19T14:22:08.512Z",
  "properties": {
    "cart_total": 49.99,
    "item_count": 3,
    "currency": "GBP"
  }
}
```

A request will then be made to the experimentation service via gRPC to fetch the definition of that event schema via the `event_key` property. This result with then be cached within an internal Caffeine Cache, and the incoming event will then have it's schema validated.

This presents issues later down the line where the data cooking service or stats service may expect a property to exist on an event, but due to a misconfiguration in th emitter, this property is not present.

Once validated, events will be transformed into the json format below, where each property is enriched with schema information. It will then be emitted to the `events` Kafka topic to be joined with experiment data in the data-cooking-service.

```json
{
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

## Validation Rules

For all incoming events, the below rules are enforced to ensure downstream consistency. Note the special case for Experiment key: Experiment exposure event must contain an experiment key, as we can't retrospectively join experiment exposure events as alone, they lack the required context.

| Rule                                                                              | Applies to                                       | On failure                                                                                              |
| --------------------------------------------------------------------------------- | ------------------------------------------------ | ------------------------------------------------------------------------------------------------------- |
| `event_key` must be present and non-empty                                         | all events                                       | 4xx — `EventIngestionException` returned to client                                                      |
| `user_details.id` must be present and non-empty                                   | all events                                       | 4xx - `EventIngestionException` returned to client                                                      |
| `sent_at` must be present                                                         | all events                                       | 4xx - `EventIngestionException` returned to client                                                      |
| `experiment_key` must be present and non-empty                                    | events with `event_key == "experiment_exposure"` | 4xx — `EventIngestionException` returned to client                                                      |
| `event_key` must resolve to a known event type in the events catalog              | all events                                       | 4xx - `EventIngestionException` returned to client                                                      |
| Every property defined on the catalog `EventType` must be present in `properties` | all events                                       | 4xx - `EventIngestionException` returned to client (Eventually this will be surfaced in the portal too) |

## API Specification

[View Events Service API Documentation →](/api/events-service)
