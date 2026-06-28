# FR-9 — Assignment

**User Story / Rationale**

As an experiment owner, I want a low latency api to query with a user id that will respond with a map of `experiment_key:variant_key` detailing the users current experiment assignments, so that the correct variants can be rendered.

**Description**

**Input**

**Output**

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

# FR-9.1 — Deterministic Experiment & Variant Assignment

**User Story / Rationale**

As an experiment owner, I want the user to see the same variant each time, so that the user experience is not dis-jointed, and so that the experiment remains statistically valid

**Description**

The implemented system must have a uniform distribution of outputs in order to ensure assignments are well distributed.

**Input**

- User ID (Assumed as string, can be int cast to string)

**Output**

- The same experiments and variants assigned to a user id per request

**Preconditions**

- The experiment is in an active serving state of A/A or A/B

**Post-Conditions**

**Error Handling**

**Security Measures**

# FR-9.2 — Cache Variant Assignment

**User Story / Rationale**

As a systems architect, I want assignments to be cached in memory so the database is protected from high volume read traffic and latency is minimised for the end user so that the render path is fast.

**Description**

Making a request to the database for every assignment request could potentially cause a system bottleneck. Although it will be functional at the MVP scale, this functional requirement aims to improve performance at a production level by caching assignment results at the assignment service, preventing the need for requests further into the system.

**Input**

**Output**

**Preconditions**

**Post-Conditions**

- Assignments are served with lower latency.

**Error Handling**

When encountering a cache miss, the system must fall back to the database and then write through the result. A decision needs to be made here on whether the cache should be written to in the event there are no assignments, or should the request continue to fall through.

**Security Measures**

# FR-11 — Client-Side React Assignment Evaluator

**User Story / Rationale**

As an experiment owner, I want a standardised way of evaluating assignments, so that users can be conditionally rendered variants in a react web app.

**Description**

Without this mechanism, prism would be limited to server-side experiments.

**Input**

- User Id

**Output**

- The user is conditionally rendered react components based on the assignments retrieved from a call to the assignment service

**Preconditions**

**Post-Conditions**

**Error Handling**

- In the event the api call fails, the user should be shown a control variant as a fallback.

**Security Measures**

- The user Id MUST be validated to ensure the user is who they say they are by the organisations proxying service. 