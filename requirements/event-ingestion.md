# FR-2 — Event Ingestion

**User Story / Rationale**

As an experiment owners I want to ingest events into prism, so that, the events can be aggregated into metrics.

**Description**

**Input**

**Output**

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

## FR-2.1 — Client Side Instrumentation

**User Story / Rationale:**  
As a data engineer, I want to be able to instrument client-side actions so that I can use these events to create metrics.

**Description:**  
Some events can only be collected on the clientside. Think hover events. These events need to be reliably instrumented. This necessitates a client side SDK to allow the instrumentation of events.

**Inputs:**

- Reusable and lightweight methods that can be called by JS on the clientside to send the event to the events service.

**Outputs:**

- The event is published to the client-side via HTTP

**Pre-Conditions:**

- The event definition must be present in the events catalog.

**Post-Conditions:**

- After a 200 response is received, the same instance of the same event must not be published again. Exactly once guarantees are required.

**Error Handling:**

- If a 500 series response code is recieved, indicating the server did something wrong, then the client should retry sending the event up to X times (Where X is a configurable property).

**Security Measures:**

- Client-side instrumentation should NOT trust user supplied identity blindly. Events sourced directly from the browser must first be validated by the organisations backend. This will prevent forged user ids from being used to target experiment results.
-

# FR-2.2 — Server Side instrumentation.

**User Story / Rationale:**  
As a data engineer, I want to be able to instrument server-side actions so that I can use these events to create metrics.

**Description:**  
Some events can only be collected Serverside. For Example payment info, financial records, integrations with other systems e.g. minutes listened.

**Inputs:**
- An event that matches the defined schema.

**Outputs:**

- The event is published via HTTP

**Pre-Conditions:**

- The event definition must be present in the events catalog.

**Post-Conditions:**

- After a 200 response is received, the same instance of the same event must not be published again. Exactly once guarantees are required.

**Error Handling:**

**Security Measures:**

- This model of ingestion requires that the user_id has been enriched on the serverside using the customers auth system. Prism will treat any user_id as valid. It is up to the customer to validate that the user_id belongs to the user who sent the event.

# FR-2.3 — Validate Events Against Catalogue

**User Story / Rationale:**  
As a data engineer, I want to only allow events to enter the system if they are defined in the events catalogue, and contain the correct event fields and property types.

**Description:**
This requirement aims to guard against human error which has the potential of impacting experiment trustworthiness. For example consider the scenario where a server-side instrumentation is incorrectly omitting the order_total field on a purchase event. The sender should be informed of this misconfiguration.

A potential extension of this requirement is to consider it's place within the wider need for an observability solution that instruments Prism itself.

**Inputs:**

- An event represented in a JSON Format in the request payload.

**Outputs:**

- If the event is malformed in that it is missing fields the sender must be informed in the response so they know to correct the payload.

**Pre-Conditions:**

- The event key contained within the request must match with an existing event key as per the catalogue.

**Post-Conditions:**

**Error Handling:**

**Security Measures:**

# FR-2.4 — Cook and Enrich Events

**User Story / Rationale:**  
As a Data Engineer, I want raw events that lack context to be enriched with user specific assignment information, so that experiment analysis does not have to join experiments with assignment data at query time.

**Description:**
This requirement is derrived from Gupta _et al_.'s (2018) discussion of data cooking.

**Inputs:**

- An event represented in a JSON Format.
- The user Id of the user who performed the event
- The Bucket the user is assigned to
- The experiments the bucket is assigned to
- The details of the experiments the experiment is assigned to
- The variants the user is assigned to within the experiment

**Outputs:**

- A fan out of events where each event is enriched for a specific experiment. E.g. if the purchase event is used in multiple experiments, then an enriched event should be created for each experiment containing the variant the user was assigned to as well as an experiment key.

**Pre-Conditions:**

- The event key contained within the request must match with an existing event key as per the catalogue.

**Post-Conditions:**

- The event will be ready to be written to the events database.

**Error Handling:**

**Security Measures:**

# FR-2.5 — Persist Cooked Events

**User Story / Rationale:**  
As a Data Engineer, I want enriched events to be persisted into a database in a structured manner, so that when calculating metric values for experiments queries can easily be created.

**Description:**

**Inputs:**

- An enriched event

**Outputs:**

- A row written to a database.

**Pre-Conditions:**

- The event has been previously enriched with experiment context

**Post-Conditions:**

- The cooked event will be persisted and available for use in metric queries

**Error Handling:**

**Security Measures:**
