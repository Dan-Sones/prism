# FR-3 — Metrics Catalogue

**User Story / Rationale**

As an experiment owner, I want to be able to create metrics within a catalogue, so that I can re-use them between experiments.

**Description**

Metrics are the basis of experiments. An interface is needed to allow common metric types to be defined and stored in the database. This will allow for re-use between experiments, instead of having to redeclare experiments each time.

**Input**

**Output**

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

# FR-3.1 Create Binary Metric

**User Story / Rationale**

As an experiment owner, I want to be able to create a binary metric, so that I can measure conversion rates.

**Description**
Binary metrics allow us to easily determine whether or not the thing being tested is successful as they show the proportion of user who matched our success condition.

**Input**

- Metric name.
- Metric identifier (key).
- Metric description.
- Numerator with aggregator defined.
- Denominator with aggregator defined.

**Output**

- The binary metric is persisted in the database

**Preconditions**

- A metric with the same key does not already exist

**Post-Conditions**

- The metric can be used in experiments.

**Error Handling**

**Security Measures**

# FR-3.2 Create Continuous Metric

**User Story / Rationale**

As an experiment owner, I want to be able to create a continuous metrics, so I can measure metrics that take place over a continuous range.

**Description**
Continuous metrics allow us to measure things over a range e.g. minutes listened.

**Input**

- Metric name.
- Metric identifier (key).
- Metric description.
- Numerator with aggregator defined.
- Denominator with aggregator defined.

**Output**

- The metric is persisted in the database

**Preconditions**

- A metric with the same key does not already exist

**Post-Conditions**

- The metric can be used in experiments.

**Error Handling**

**Security Measures**

# FR-3.3 Create Ratio Metric

**User Story / Rationale**

As an experiment owner, I want to be able to create a ratio metrics, so I can express the relationship between two values.

**Description**
Ratio metrics allow the measurement of events relative to other types of events. This comes with considerably deeper implementation complexity.

**Input**

- Metric name.
- Metric identifier (key).
- Metric description.
- Numerator with aggregator defined.
- Denominator with aggregator defined.

**Output**

- The metric is persisted in the database

**Preconditions**

- A metric with the same key does not already exist

**Post-Conditions**

- The metric can be used in experiments.

**Error Handling**

**Security Measures**
