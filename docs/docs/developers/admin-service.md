---
sidebar_position: 2
---

# Admin Service

The Admin Service is the core of the Prism.

See the following TOC for a breakdown of the Admin Service's responsibilities and capabilities:

- [Event Catalog Management](#event-catalog-management)
- [Metric Catalog Management](#metric-catalog-management)
- [Experiment Configuration Management](#experiment-configuration-management)
- [API Specification](#api-specification)
- [GRPC Specification](#grpc-specification)

## API Specification

[View Admin Service API Documentation →](/api/admin-service)

## GRPC Specification

Please see `prism/services/proto/assignment` for the GRPC specification for this service.

## Event Catalog Management

Prism allows you to define a catalog of events to allow the architecture to store and consume them.

Events can be ingested in one of two ways:

1. Through [Open Feature](./open-feature.md)'s `.track()` functionality, allowing events to be tracked on the client-side.
2. Through manual event delivery to the Events service (Useful for backend events, e.g. payment completed, email sent etc).

Events defined in the event catalog can then be used in metric definitions, allowing prism to calculate results on variables more complex than simple click counts.

Events exist in the following format:

```json
{
  "id": "4522d155-df00-4877-b77a-9ce999c99446",
  "name": "purchase_completed",
  "description": "Fired when a user completes a purchase",
  "createdAt": "2026-02-18T21:53:44.867961Z",
  "fields": [
    {
      "id": "ad8644e5-9b11-4d95-9b76-e46b22458f94",
      "name": "Currency",
      "fieldKey": "currency",
      "dataType": "string"
    },
    {
      "id": "a56852d8-1c84-450b-b409-41a2a8ecf823",
      "name": "Order Total",
      "fieldKey": "order_total",
      "dataType": "float"
    }
  ]
}
```

### Event Type

| Field         | Type              | Description                                                    |
| ------------- | ----------------- | -------------------------------------------------------------- |
| `id`          | string (UUID)     | Unique identifier for the event                                |
| `name`        | string            | Display name of the event                                      |
| `description` | string            | Human-readable description of what the event represents        |
| `createdAt`   | string (ISO 8601) | Timestamp when the event was created                           |
| `fields`      | array             | Array of field definitions that can be tracked with this event |

### Event Field Properties

| Field      | Type          | Description                                                                |
| ---------- | ------------- | -------------------------------------------------------------------------- |
| `id`       | string (UUID) | Unique identifier for the field                                            |
| `name`     | string        | Display name of the field (to be used in the ui)                           |
| `fieldKey` | string        | The key used when sending this field in events                             |
| `dataType` | string        | The data type of the field (e.g., `string`, `float`, `integer`, `boolean`) |

Think of event fields as the keys in the event payload. See the below JSON. As a user you may want to track payment events and extract the `order_total` field from those events to use in your metrics. In the case of nested events, the `fieldKey` also serves as the path to the field in the event payload using dot-notation. Therefore, if we wanted to use the `order_total` property in our metric calculations the `fieldKey` for that event field within an EventType would be `purchase_details.order_total`.

```json
{
  "event": "purchase_completed",
  "properties": {
    "customer": {
      "id": "123",
      "loyalty_member": true
    },
    "purchase_details": {
      "order_total": 100.0,
      "currency": "USD"
    }
  }
}
```

## Metric Catalog Management

## Experiment Configuration Management
