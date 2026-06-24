---
sidebar_position: 2
---

# Experimentation Service

The Experimentation Service is the core of the Prism. It is responsible for the definition and management of Events, Metrics and Experiments. It interfaces with postgres to store definitions, and makes requests out to the analysis service to perform statistical operations.

See the following TOC for a breakdown of the Experimentation Service's responsibilities and capabilities:

- [Event Catalog Management](#event-catalog-management)
- [Metric Catalog Management](#metric-catalog-management)
- [Experiment Configuration Management](#experiment-configuration-management)
- [API Specification](#api-specification)
- [GRPC Specification](#grpc-specification)

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
| `event_key`   | string            | The unique key used to identify the event                      |
| `description` | string            | Human-readable description of what the event represents        |
| `created_at`  | string (ISO 8601) | Timestamp when the event was created                           |
| `fields`      | array             | Array of field definitions that can be tracked with this event |

### Event Field Properties

Event Field Properties represent one of the values in your Event. The context provided here determines how your item is stored in the database and therefore it's use in metric's, so it is important that event fields are correctly typed.

| Field       | Type          | Description                                                                |
| ----------- | ------------- | -------------------------------------------------------------------------- |
| `id`        | string (UUID) | Unique identifier for the field                                            |
| `name`      | string        | Display name of the field (to be used in the ui)                           |
| `field_key` | string        | The key used when sending this field in events                             |
| `data_type` | string        | The data type of the field (e.g., `string`, `float`, `integer`, `boolean`) |

## Metric Catalog Management

A lot of thought has gone into the design of the metrics catalog. By defining metrics as re-usable entities, their use in experiments is much more straight forward. Define a metric once and you can use it in as many experiments as you want.

Metrics are defined using a multi-record database system. The parent record contains information about the metric, such as its name, an identifiable key, description, its type, and analysis unit. Metric components can be thought of as components of a fraction; you have a numerator and a denominator. Numerators and denominators both point at events and event fields defined within the above-mentioned event catalogue. There is then a number of predefined aggregations that can be applied to event fields in order to turn values into interpretable metrics.

This relationship can be visualised in the below ERD.

![ERD](/img/erd.png)

The below JSON shows an example metric definition for a binary conversion rate metric.

Although Prism currently only supports binary success metrics the database design and user interface used to create metrics is flexible enough to allow for users to easily define continuous and ratio metrics in the future.

```json
{
  "id": "0b43d53e-c722-4222-964d-04e24af872db",
  "name": "Purchase Conversion Rate",
  "metric_key": "purchase_conv_rate",
  "description": "The percentage of exposed users who made at\nleast one purchase.",
  "created_at": "2026-05-23T09:10:34.458352+01:00",
  "metric_type": "ratio",
  "analysis_unit": "user",
  "metric_components": [
    {
      "id": "7490bd70-096f-4e95-972c-505a8bc098ab",
      "role": "numerator",
      "event_type": {
        "id": "088436c4-ceab-4cd7-b09e-45756b00f043",
        "name": "purchase",
        "event_key": "purchase",
        "version": 1,
        "description": "Fires when a user completes a purchase",
        "created_at": "2026-05-23T09:10:34.426193+01:00",
        "fields": [
          {
            "id": "c2668043-047d-4e29-bca3-e8bf0b48bdc5",
            "name": "Order Total",
            "field_key": "order_total",
            "data_type": "float"
          }
        ]
      },
      "aggregation_operation": "COUNT_DISTINCT",
      "system_column_name": "user_id"
    },
    {
      "id": "ceead965-eb92-4678-947c-06c677f41f33",
      "role": "denominator",
      "event_type": {
        "id": "412a1b99-ad13-40d1-8b11-e7ef8d082447",
        "name": "Experiment Exposure",
        "event_key": "experiment_exposure",
        "version": 1,
        "description": "Fired when a user is exposed to an experiment.",
        "created_at": "2026-05-23T09:08:56.951375+01:00",
        "fields": []
      },
      "aggregation_operation": "COUNT_DISTINCT",
      "system_column_name": "user_id"
    }
  ],
  "is_binary": true
}
```

As well as being able to select fields on an event when creating a metric, users can also select system fields (e.g. system context attached to events). The most common of these is the user ID. As illustrated in the example above, the aggregation is pointed towards the system column name. This is implied from the absence of a targeted event field. The count distinct aggregation will apply to the user ID, thus allowing the conversion rate metric to be formed.

The JSON below shows how a continuous metric — "Average Order Value Per Exposed User" — _could_ be defined. The numerator targets the `order_total` event field with a `SUM` aggregation, while the denominator counts distinct users from the exposure event as before. The numerator references `aggregation_field_id` (pointing at a real event field) rather than `system_column_name`, and how `is_binary` is `false`.

Although in the final deliverable events defined in this format cannot be used, it shows the flexibility of the schema.

```json
{
  "id": "f1c34a2e-7d6b-4b81-9c2f-b3a5e6f8a1c2",
  "name": "Average Order Value Per Exposed User",
  "metric_key": "avg_order_value_per_user",
  "description": "Total order revenue divided by the number of distinct users exposed to the experiment.",
  "created_at": "2026-06-01T20:14:22.000000+01:00",
  "metric_type": "ratio",
  "analysis_unit": "user",
  "metric_components": [
    {
      "id": "3a7c91bd-4f88-4e25-9c0a-1b5e3f6c8d4f",
      "role": "numerator",
      "event_type": {
        "id": "088436c4-ceab-4cd7-b09e-45756b00f043",
        "name": "purchase",
        "event_key": "purchase",
        "version": 1,
        "description": "Fires when a user completes a purchase",
        "created_at": "2026-05-23T09:10:34.426193+01:00",
        "fields": [
          {
            "id": "c2668043-047d-4e29-bca3-e8bf0b48bdc5",
            "name": "Order Total",
            "field_key": "order_total",
            "data_type": "float"
          }
        ]
      },
      "aggregation_operation": "SUM",
      "aggregation_field_id": "c2668043-047d-4e29-bca3-e8bf0b48bdc5"
    },
    {
      "id": "9d4f2b7a-6e21-4c93-8f1d-2a5b7c9e3f0a",
      "role": "denominator",
      "event_type": {
        "id": "412a1b99-ad13-40d1-8b11-e7ef8d082447",
        "name": "Experiment Exposure",
        "event_key": "experiment_exposure",
        "version": 1,
        "description": "Fired when a user is exposed to an experiment.",
        "created_at": "2026-05-23T09:08:56.951375+01:00",
        "fields": []
      },
      "aggregation_operation": "COUNT_DISTINCT",
      "system_column_name": "user_id"
    }
  ],
  "is_binary": false
}
```

## Experiment Configuration Management

## GRPC Specification

Please see `prism/services/proto/assignment` for the GRPC specification for this service.
