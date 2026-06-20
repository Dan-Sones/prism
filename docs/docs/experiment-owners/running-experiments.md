---
sidebar_position: 4
---

# Experiments

Experiments in Prism are a multi-stage process.

1. Designing Experiments
2. Collect Baseline Variance In the A/A Phase
3. Perform Sample Size Calculations and Allocate Experiment Dates
4. Run A/B Test

# Designing Experiments

To create an experiment:

1. Navigate to the /experiments/create page where you should see the start of the journey

![Create Experiment Form #1](/img/create-experiment-form-1.png)

2. **Experiment Name** This is self explanatory, however like naming metrics it's important to communicate the experiments aims.

3. **Feature Flag Key** Extra care must be taken when inputting into this box. The feature flag key is the flag that is to be interpreted by the open feature provider used to conditionally render components to the user in the frontend. This value must be globally unique.

4. **Hypothesis** the term "hypothesis" has numerous meanings in the context of experimentation. This input refers to the experiment hypothesis, rather than a hypothesis in a hypothesis test[^confidence_a]. A good template to use for this:

   :::tip Hypothesis Template
   **Based on** [prior knowledge], **we believe that** [theory about user need]. **We think that** [doing this/building this feature/creating this experience] **for** [these people/personas] **will achieve** [these outcomes]. **We will know this is true when we see** [metric results].
   :::
   [^confidence_b]

5. **Description** The description allows you as the experiment owner to note any important details to other users who may observe your experiment.

6. **Metric Definitions** This UI allows you to define metrics. The success metric requires you to input a Minimum Detectable Effect (MDE). The MDE serves two purposes. Firstly it is used as an input to the sample size calculator after the A/A phase to ensure the sample size is able to detect the change you want to detect. The MDE also acts as the practical significance percentage increase you want to see in a treatment variant compared to control to justify shipping the feature.
   :::note MVP Scope
   In the current release, only **Success Metrics** are fully supported. Guardrail, deterioration, and quality metric types are declarable but statistical tests for these metrics are yet to be implemented.
   :::

![Create Experiment Form #2](/img/create-experiment-form-2.png)

7. **Variant Definitions** This page allows you to configure the Control and Treatment Variants. The Control Variant should represent the existing user experience and the Treatment variant should represent the new Experience.

:::warning
It is **Vital** that the variant_key fields match the keys configured within the OpenFeature components. If there is a misconfiguration here, users may be stuck seeing exclusively the control variant (Assuming it is configured as the default)
:::

8. **A/A Test** This screen allows you to start the A/A test. See #A/A

![Create Experiment Form #3](/img/create-experiment-form-3.png)

# The A/A Test

The A/A assigns the entire population and collects the same data the A/B test would. This data is then used to calculate variance, which is then used as an input to a sample size calculator.

The A/A test runs for 7 days in order to account for the day of week effect.

# Configuring The A/B Test

When the A/A test is complete you will be able to allocate a percentage of the total population to the experiment. This will randomly sample from the available buckets. You must also specify the date range of the experiment.

# Evaluating the A/B test

When the experiment end date is complete you will be presented with the following UI.

![Create Experiment Form #5](/img/create-experiment-form-5.png)

For Binary Success metrics Prism runs a Z-Test. It's results are run through a decision rule to provide a shipping recommendation based on Practical Significance (Based on your stated MDE) and Statistical Significance. The range of outputs of the decision rule is shown in the table below. This decision rule is an interpretation of Kohavi et al.'s Interpretation metric results in the book trust worthy online controlled experiments [^kohavi].

| Outcome           | Statistical Significance | Practical Significance | Recommendation       | Action                                                                                                                                     |
| ----------------- | ------------------------ | ---------------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| Clear Winner      | ✅                       | ✅                     | **Recommend**        | Ship the treatment variant. It met your MDE and was statistically significant                                                              |
| No Effect         | ❌                       | ❌                     | **Do Not Recommend** | Do not ship. Results were not Practically OR significantly significant. Rethink your hypothesis and maybe re-run with a larger sample size |
| Negligible Effect | ✅                       | ❌                     | **Do Not Recommend** | The effect is real but does not meet you MDE so can't be recommended                                                                       |
| Regression        | ✅                       | ❌                     | **Do Not Recommend** | Treatment caused a statistically significant regression vs control. Do not ship.                                                           |
| Borderline        | ✅                       | ⚠️                     | **Inconclusive**     | Re-run the experiment with a larger sample size                                                                                            |

[^confidence_a]: Spotify Confidence (2026) _Lesson 1: Introduction to hypothesis tests_ Available At: https://confidence.spotify.com/bootcamp/hypothesis-testing/what-is-hypothesis-testing (Accessed: 28 May 2026)

[^confidence_b]: Spotify Confidence (2026) **Lesson 2: Experiment hypothesis** Available At: https://confidence.spotify.com/bootcamp/intro-course/experiment-hypothesis (Accessed: 28 May 2026)

[^kohavi]: Kohavi, R., Tang, D. and Xu, Y. (2020) Trustworthy Online Controlled Experiments: A Practical Guide To A/B Testing. Cambridge University Press.
