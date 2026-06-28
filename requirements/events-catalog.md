# FR-1 — Events Catalogue

## User Story / Rationale

As an experiment owner, I want the system to maintain a known catalogue of events so that I can create new metrics.

**Description**

The Events Catalogue shall act as the source of truth for event ingestion and metric definitions.
The catalogue enables the Metrics Catalogue and experiment functionality.

**Important:** This FR does NOT include update and delete options for event. For the case of the MVP Events are considered immutable. This is due to the dependency chain from _event -> metric -> experiment_. Implementing this would include either outright denying the user from deleting the event while it forms the start of any number of dependency chains, or putting it in a queue for deletion when it next becomes available.

**Input**

**Output**

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

## FR-1.1 — Create Event

**User Story / Rationale**

As an experiment owner, I want to define the structure and properties of an event so that it can be stored in the Events Catalogue and used to build metrics.

**Description**

To build a metric, events need to be defined. This FR encompasses the creation of Event Definitions within the events catalog.

**Input**

- Event Name - A string used in the UI when referencing events
- Event Key - A unique string system identifier used to classify event request bodies.
- Description - A string description of the event.
- Event Fields - An array of objects:
  - Field Name - A string used in the ui when referencing an event field
  - Field Key - A locally unique string key
  - Data Type - An ENUM of _String, Int, Float, Boolean, Timestamp_ representing the data type of the above field.

**Output**

- Confirmation to the user the event was successfully implemented
- The event schema is persisted in the database.

**Preconditions**

- The users desired event key must not already be in use.

**Post-Conditions**

- The event and its fields are available to build metrics within the metrics catalog.
- The events service will accept client-side and server side instrumentation events matching the defined schema.

**Error Handling**

- The event key MUST be unique

**Security Measures**

- The event key MUST be validated by a regex to ensure that it only contains letters, numbers, and underscores. This will ensure it remains interpretable across the system as well as lowering the risk of sql injection. This should be validated on the client side AND serverside

## FR-1.2 - View Event Details and Usage Graphs

**User Story / Rationale**

As an experiment owner, I want visibility on events flowing into the system, so that I know I have instrumented the system correctly.

**Description**

Without some visbility into this, experiment owners are left blind as to whether events are actually flowing into their system unless they manually review system logs, kafka topics, or the database

**Input**

- Event Key - A unique string system identifier

**Output**

- A time series graph visualising the ingest of events into the system across a number of selectable timescales.

**Preconditions**

- The event must be configured.

**Post-Conditions**

**Error Handling**

**Security Measures**
