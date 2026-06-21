---
sidebar_position: 3
---

# Metrics Catalog

Metrics transform raw user behavior (events) into actionable insights.[^confidence_a]

# Supported Metrics

In the MVP state that this delivery represents, Prism supports binary metrics which measure did or did not something happen [^confidence_a]. Binary Metrics are powerful as they represent proportions which we often want to improve for example:

1. Conversion Rates - _Did the user purchase something when exposed to the new Checkout UI?_
2. Sign Up Rates - _Did the user register for an account when exposed to the register pop-up?_
3. Click-Through Rate - _Did the user engage with this suggested article more than this one?_
4. Unsubscribe Rate - _Did the user like the new style of marketing email or did it cause them to unsubscribe?_

This list could carry on forever which is precisely the reason that metric development was scoped here. They are perhaps the easiest to interpret as well as being extremely powerful with what they tell us.

# Metric Catalog

Metrics are defined at the platform level to allow metric definitions to be re-used between experiments.

## Creating a New Metric

Before you can create a metric it is important to first have the event type that forms the parts of your metrics defined in the events catalog, see [Events Catalog Documentation](./events-catalog.md).

Binary Metrics measure the proportion of users that 'did something' given a population of exposed users exposed to the experiment.

It is helpful to think of binary metrics as a fraction where the denominator of the fraction is those users who were exposed, and the numerator is of those users who did something. Where the output is a proportion between 0 and 1 representing the percentage of exposed users who converted.

:::caution
One **IMPORTANT** thing to consider with binary metric is that it best deals with **unique users** both in the population (Denominator) and in the proportion (Numerator). Take a purchase conversion rate metric, a user can in theory purchase something multiple times (Numerator), as can they be exposed to a new checkout UI multiple times (denominator). This does not mean however that we should count this individually, they are the same user after all. If we were to remove the DISTINCT from the COUNT_DISTINCT that is applied on binary metrics, a user that visits the checkout numerous times but only buys once would negatively affect the metric despite there being a positive outcome.
:::

1. Visit the /metrics-catalog/create page. Where you should see the following form.

![Create Metric Form](/img/create-metric-form.png)

2. **Metric Name** Choose a name that easily and accurately identifies your metric so it is easy to re-use between experiments. Metric names should be descriptive and clarify what you are measuring [^confidence_b].

3. **Metric Key** The metric is a **UNIQUE** system identifier for metrics.

4. **Metric Type** and **is Binary** determines the form rendered below. In the current deliverable, this is just Binary Metrics (A subset of ratio metrics) so always set these to "Ratio" and tick is binary.

5. You will now be presented with a numerator and denominator, this goes back to the fraction analogy described above. To correctly implement a binary metric.

6. For the Numerator: Search for the event key of which proportion you want to measure, then select User ID as the field, and COUNT_DISTINCT as the aggregation operation

7. For the Denominator: Search for the experiment_exposure event and again select User ID as the field, and COUNT_DISTINCT as the aggregation operation.

8. You will now be able to create the metric and use it in experiments!

[^confidence_a]: Spotify Confidence (2026) _Lesson 1: What is a metric?_. Available At: https://confidence.spotify.com/bootcamp/intro-to-metrics/what-is-a-metric (Accessed: 28 May 2026)

[^confidence_b]: Spotify Confidence (2026) _Lesson 6: Interpretability_. Available At: https://confidence.spotify.com/bootcamp/intro-to-metrics/interpretability (Accessed: 28 May 2026)
