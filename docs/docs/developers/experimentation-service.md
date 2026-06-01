---
sidebar_position: 2
---

# Experimentation Service

The Experimentation Service is the core of the Prism.

See the following TOC for a breakdown of the Experimentation Service's responsibilities and capabilities:

- [Event Catalog Management](#event-catalog-management)
- [Metric Catalog Management](#metric-catalog-management)
- [Experiment Configuration Management](#experiment-configuration-management)
- [API Specification](#api-specification)
- [GRPC Specification](#grpc-specification)

## API Specification

[View Experimentation Service API Documentation →](/api/experimentation-service)

## GRPC Specification

Please see `prism/services/proto/assignment` for the GRPC specification for this service.

## Event Catalog Management

Prism allows you to define a catalog of events to allow the architecture to store and consume them.

Events can be ingested through an api call to the events service.

Events defined in the event catalog can then be used in metric definitions, allowing prism to calculate results on variables more complex than simple click counts.

Events exist in the following format:

```json
{
  "id": "04d83743-3a17-40d1-bccb-76b56f8987bb",
  "name": "Purchase Successful",
  "event_key": "purchase_successful",
  "description": "The purchase event is triggered when a user successfully uses their method of payment to complete an order. The event contains a number of properties about the purchase.",
  "created_at": "2026-06-01T19:53:58.172026+01:00",
  "fields": [
    {
      "id": "e7484700-c861-449c-ac96-f06ef2ebc326",
      "name": "Basket Item Count",
      "field_key": "basket_item_count",
      "data_type": "int"
    },
    {
      "id": "e5282870-737f-499a-b14f-88791342d933",
      "name": "Currency",
      "field_key": "currency",
      "data_type": "string"
    },
    {
      "id": "1f4b5354-985c-4a7b-8fce-6e7b1cd83135",
      "name": "Is Loyalty Member",
      "field_key": "is_loyalty_member",
      "data_type": "string"
    },
    {
      "id": "d0b980ea-7563-461a-9879-bc6cf3c4e4e4",
      "name": "Order Time",
      "field_key": "order_time",
      "data_type": "timestamp"
    },
    {
      "id": "05628ab6-9dec-459a-b661-180a11de83dc",
      "name": "Order Total",
      "field_key": "order_total",
      "data_type": "float"
    }
  ]
}
```

### Event Type

| Field         | Type              | Description                                                    |
| ------------- | ----------------- | -------------------------------------------------------------- |
| `id`          | string (UUID)     | Unique identifier for the event                                |
| `name`        | string            | Display name of the event                                      |
| `eventKey`    | string            | The unique key used to identify the event                      |
| `description` | string            | Human-readable description of what the event represents        |
| `createdAt`   | string (ISO 8601) | Timestamp when the event was created                           |
| `fields`      | array             | Array of field definitions that can be tracked with this event |

### Event Field Properties

Event Field Properties represent one of the values in your Event. The context provided here determines how your item is stored in the database and therefore it's use in metric's, so it is important that event fields are correctly typed.

| Field      | Type          | Description                                                                |
| ---------- | ------------- | -------------------------------------------------------------------------- |
| `id`       | string (UUID) | Unique identifier for the field                                            |
| `name`     | string        | Display name of the field (to be used in the ui)                           |
| `fieldKey` | string        | The key used when sending this field in events                             |
| `dataType` | string        | The data type of the field (e.g., `string`, `float`, `integer`, `boolean`) |

## Metric Catalog Management

## Experiment Configuration Management
