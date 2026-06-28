# FR-4 — Create Experiments

**User Story / Rationale**

As an experiment owner, I want to be able to run experiments to measure metric values between a control and treatment variant, so that I can make informed business decisions.

**Description**

**Input**

- Experiment name.
- Feature flag identifier.
- Hypothesis.
- Experiment description.
- One or more metrics selected from the metrics catalogue (See sub-requirements for metric roles)
- Expected movement direction for each selected metric.
- Control variant name.
- Control variant identifier.
- Treatment variant name.
- Treatment variant identifier.

**Output**

- The experiment exists within the database.

**Preconditions**

- The experiment key / feature flag key is not already in use

**Post-Conditions**

- The experiment is associated with it's key / feature flag and will be served by the assignment service when the experiment begins.
- The experiment can enter an A/A test phase.

**Error Handling**

**Security Measures**

# FR-4.1 Success Metrics

**User Story / Rationale**

As an experiment owner, I want to associate metrics to experiments as success metrics, so that I can measure the success case of my experiments.

**Description**

**Input**

- A metric that exists within the metrics catalogue.
- Minimum detectable effect.

**Output**

**Preconditions**

**Post-Conditions**

- The metric is attached to the experiment as a success metric
- The metric is interpreted by checking if the experiment achieved it's intended outcome as per the hypothesis.

**Error Handling**

**Security Measures**

# FR-4.2 Guardrail Metrics

**User Story / Rationale**

As an experiment owner, I want to associate metrics to experiments as guardrail metrics, so that I can be assured experiments that I run do not negatively affect metrics by a certain margin.

**Description**

**Input**

- A metric that exists within the metrics catalogue.
- Non Inferiority margin.

**Output**

**Preconditions**

**Post-Conditions**

- The metric is attached to the experiment as a guardrail metric
- The metric is interpreted by the ensuring the non-inferiority margin has not been exceeded.

**Error Handling**

**Security Measures**

# FR-4.3 Deterioration Metric

**User Story / Rationale**

As an experiment owner, I want to associate metrics to experiments as deterioration metrics, so that I can be assured experiments that I run do not deteriorate as a result of the experiment.

**Description**

**Input**

- A metric that exists within the metrics catalogue.
- Non Inferiority margin.

**Output**

**Preconditions**

**Post-Conditions**

- The metric is attached to the experiment as a Deterioration metric
- The metric is interpreted by the ensuring the non-inferiority margin has not been exceeded.

**Error Handling**

**Security Measures**

# FR-4.4 Quality Metrics (Sample Ratio Mismatch)

**User Story / Rationale**

As an experiment owner, I want to associate measurable values (Like SRM) to my experiment, so that i can detect issues with experiment trustworthiness.

**Description**

This FR is potentially misleading in it's apparent lack of complexity. In reality here there is need for some serious design decisions. For example, should SRM be a metric that sits within the metrics catalogue, or should it be something separate. The answer to this will influence the way the metrics catalogue itself it designed.

The remainder of this FR is left empty as it's not apparent. I need to do more research into this area as an SRM implementation cannot be half baked.

**Input**

**Output**

- The quality check is attached to the experiment.
- The quality check is used to detect issues with experiment setup, assignment, exposure, or data collection.

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

# FR-5 A/A Test (Sample Size Calculator)

**User Story / Rationale**

As an experiment owner, I want to complete an A/A test, so that I can establish the baseline conversion rate of a metric.

**Description**
The baseline conversion rate of a metric, and inherently it's variance determines the sample size that is required in order to reach statistical power. The A/A implementation will provide a mechanism in which to measure this baseline to then feed to a sample size calculator.

**Input**

- 7 Days worth of data (Enforced in order to account for day of week effect)
- MDE / NIM from associated metrics

**Output**

- A required sample size presented in the UI.

**Preconditions**

- 100% of the assigned population must be served the control variant for the duration of the A/A phase.

**Post-Conditions**

- The experiment owner can use the calculated sample size to inform them what percentage of the population to assign to the experiment.

**Error Handling**

**Security Measures**

# FR-6 A/B Test (Cumulative Metrics)

**User Story / Rationale**

As an experiment owner, I want the system to calculate whether the treatment group changes were statistically significant at the end of the experiment, so that I can get a clear "Ship / No Ship" recommendation.

**Description**

**Input**

- Aggregated cooked event data for both Control and Treatment Variants.
- The metrics associated with the experiment
- The MDE / NIM associated with the experiment.

**Output**

- (Exact properties are dependant on statistical test used)
- A clear ship / no ship recommendation
- Whether or not the required sample size was met.

**Preconditions**

- The experiment must have completed it's pre-defined duration

**Post-Conditions**

- The experiment owner is able to act on a statistically sound ship / no ship recommendation

**Error Handling**

- Clearly display when the required sample size was not met, indicating the experiment was unable to reach statistical power.

**Security Measures**

# FR-6.1 A/B Test (Window-Based Metrics)

**User Story / Rationale**

As an experiment owner, I want my experiment to use window-based metrics so that I can consider time and novelty effect in my metric analysis, rather than cumulatively using all the data for a user since their exposure.

**Description**

**Input**

**Output**

**Preconditions**

**Post-Conditions**

**Error Handling**

**Security Measures**

# FR-7 Cancel Experiment

**User Story / Rationale**
As an experiment owner, I want to immediately cancel an active experiment and revert all traffic to the control variant, so that I can mitigate the impacts of faulty treatments.

**Description**

**Input**

- Experiment identifier

**Output**

- The experiment state is updated to cancelled
- The assignment service "forgets" about the experiment, so it is no longer served to assignment requests.

**Preconditions**

- The experiment has not already been cancelled

**Post-Conditions**

- Users only see the control variant

**Error Handling**

**Security Measures**

# FR-8 Experiment Time-Deterministic State Management

**User Story / Rationale**
As an experiment owner, I want the system to automatically manage the state transitions of my experiment based on the durations I have inputted, so that traffic allocation happens automatically without manual intervention.

**Description**

**Input**

- A/A Phase Start Time (Automatically scheduled to the next day's UTC midnight).
- A/A Phase Duration (Fixed to 7 days to account for day-of-week seasonality).
- A/B Phase Start Time & End Time (User-defined).

**Output**

- The experiment state is derived from the current time and the configured A/A and A/B phase windows.
- During the A/A phase, assigned users receive the control variant.
- During the A/B phase, assigned users receive control or treatment variants according to the configured variant bounds.

**Preconditions**

- To transition into the A/B Phase, the A/A phase must be complete

**Post-Conditions**

- The assignment service serves traffic based on the real time state of an experiment
- The UI displays the current state of the experiment
- Completed or cancelled experiments are excluded from active assignment responses.

**Error Handling**

**Security Measures**

# FR-10 — Exclusivity Groups

**User Story / Rationale**

As an experiment owner, I want to be able to assign experiemnts into "exclusivity groups" so that there is a reduced risk of interaction effects.

**Description**

Sometimes more than one experiment might be taking on one local piece of an application at a time. Users cannot be assigned to both of these experiments at the same time, otherwise there is a risk of the experiment being simply unfeasible in that a button might need to be two colours at once as per an assignment response. Exclusivity groups exist as a mechanism to prevent buckets from being assigned to more than one experiment within an exclusivity group. For example, you may declare an exclusivity group for the home page to ensure the user is only engaged in one home page experiment. This will also prevent interaction effects.

**Input**

- An exclusivity group name

**Output**

- The exclusivity group is persisted and can be associated with experiments.

**Preconditions**

**Post-Conditions**

- Another experiment within the exclusivity group cannot be assigned the same bucket/s.
- When an experiment is complete, the bucket is automatically removed from the exclusivity group, freeing it up.

**Error Handling**

**Security Measures**
