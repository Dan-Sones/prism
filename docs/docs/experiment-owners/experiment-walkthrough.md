# Experiment Walkthrough

This piece of documentation serves as a walkthrough, teaching you how to instrument your own server-side API and get assignments working on your React web app.

## Experiment Context

This walkthrough acts on a fictional company, Power Motors, from The Simpsons Season 2, Episode 15 "Oh Brother, Where Art Thou?".

The designer of the Power Motors Homer product page (seen below) believes that changing the Buy Now button from black to red will improve the purchase conversion rate of the page by 2%.

A user is considered exposed when the Power Motors Analytics API receives a page viewed event. This is triggered when the user loads onto the product page. For the sake of the experiment and this walkthrough, a purchase is considered made when the user presses the Buy Now button. (In reality, you would actually place the equivalent event ingestion request after the user has had their payment processed rather than on the Buy Now button. However, this is semantics. In the case of this example).

## Walkthrough

### Prism Docker Deployment

This tutorial aims to mirror deploying Prism in a production environment as closely as possible. Therefore, we will be deploying Prism within its native Docker Compose environment.

Please see the docker-compose file at: `examples/powell-motors/docker-compose.yml`. It is pre-configured with all of the required variables and communication between services. The first thing you must do is run the following command to copy the env example.

```
cp .env.example .env
```

You shouldn't need to make any changes in here for the case of the demo; however, feel free to take this file and use it for your own deployment.

Next, to run Prism, run the following command: `make run-prism`. This will pull in all of the docker images for each of Prism's services, as well as images for Postgres, Redis, Clickhouse, etc. And then run them on your machine using the environment variables from the .env file.

If all has gone well, you should be able to access the portal on localhost port 3000.

Next, in the terminal, run `make run-powell`. This will start the Power Motors website and Spring Boot backend. Accessible on port 5147 and 8090, respectively.

### Creating Events

In this scenario, we care about two different event types:

- Experiment exposure events
- Purchase events

Experiment exposure event schemas are already created for us by the fly-wide database migrations as you boot Prism. For purchase events, however, these are an event specific to our domain, so we need to create them within the events catalog.

**1.** Open the portal and visit the events catalog.

**2.** Click the Create Event button and fill in the following details:

| Field     | Value    |
| --------- | -------- |
| Name      | Purchase |
| Event Key | purchase |

The event should contain the following fields:

| Field      | Type   |
| ---------- | ------ |
| product_id | string |
| amount_gbp | float  |
| order_id   | string |

**3.** Press Create Event Type

### Creating Metrics

As detailed in experiment context, The designer believes that the purchase conversion rate metric will be improved. The purchase conversion metric is of users who have visited the page and have purchased the homer.

This is a binary metric which we can create in our metrics catalogue.

**1.** Visit the metrics catalog

**2.** Press the Create Metric button and input the following fields:

| Field       | Value                    |
| ----------- | ------------------------ |
| Metric Name | Purchase Conversion Rate |
| Metric id   | purchase_conv_rate       |
| Metric Type | Ratio                    |
| is Binary   | True                     |

**3.** You will now be presented with a Numerator and Denominator Field Box.

Think of this as a fraction, as described in the rest of the documentation. On the top of the fraction, the numerator, we want to know the number of users who have converted (e.g., they have bought a homer). Therefore, we must input the Event Type key for a purchase event. We are concerned with the number of unique users who have purchased a home. Therefore, once we have selected the purchase event type key, we should select the user id field and the count distinct numerator aggregation.

| Field                 | Value          |
| --------------------- | -------------- |
| Event Type Key        | purchase       |
| Event Field Key       | User ID        |
| Numerator Aggregation | COUNT_DISTINCT |

For the denominator, we want to know how many users have visited the homer product page. To do this, the denominator should be experiment exposure events. Don't worry about how this is linked to an experiment. This is handled automatically in the background by the data cooking service. Input the following fields:

| Field                 | Value               |
| --------------------- | ------------------- |
| Event Type Key        | experiment_exposure |
| Event Field Key       | User ID             |
| Numerator Aggregation | COUNT_DISTINCT      |

**4.** Press Create Metric. This metric is now available to use in experiments.

### Server-Side instrumentation

Before we can create and begin our experiment, we must instrument our server to send the required events. If you can remember from the previous steps, we are concerned with two event types here:

1. An experiment exposure event which details that the user has seen the homer product page
2. A purchase event when the user "purchases" the home row by pressing the button

If you look in the EventController.java file of the provided Spring Boot project, You will see two endpoints:

- /view-product endpoint
- /purchase endpoint

The /view-product endpoint will be called when the user visits the page on the front end. Imagine that this event would be connected with your existing analytics code.

The purchase event, for the sake of this demonstration, is called when the user presses the Buy Now button. Realistically, however, you would not trigger this event until the user's payment has gone through, but for the sake of simplicity it will be called as soon as the user presses the button.

For the sake of this walkthrough, the instrumentation has been done for you. However, before beginning an experiment of your own, you must instrument your server with an experiments required events so they reach Prism.

Before the experiment begins, you must also have redeployed your backend service so that when the experiment starts, events actually flow into Prism.

### Client-Side Evaluation

The next step is configuring your frontend with our custom-developed Open Feature Provider so that a user's assignments are correctly represented. As detailed in the documentation, prism uses a bucketing model, meaning a certain number of buckets are assigned to an experiment. In an experiment, you will either be assigned to the control or treatment variant.

The control variant should be the existing experience, whereas the treatment should be the new experience.

If you view ProductDetails.tsx, See reference to the Buy Now button component. This button contains the following code.

```tsx
import { useStringFlagValue } from "@openfeature/react-sdk";
import PrimaryButton from "../../components/buttons/PrimaryButton";

type ButtonVariant = "button-black" | "button-red";

const CONTROL_VARIANT: ButtonVariant = "button-black";
const TREATMENT_VARIANT: ButtonVariant = "button-red";

const variantClassName: Record<ButtonVariant, string> = {
  "button-black": "bg-black-700!",
  "button-red": "bg-red-700!",
};

interface BuyNowButtonProps {
  handlePurchase: () => void;
}

const BuyNowButton = (props: BuyNowButtonProps) => {
  const { handlePurchase } = props;
  const evaluatedVariant = useStringFlagValue(
    "the-homer-product-page-button",
    CONTROL_VARIANT,
  );

  const buttonVariant: ButtonVariant =
    evaluatedVariant === TREATMENT_VARIANT
      ? TREATMENT_VARIANT
      : CONTROL_VARIANT;

  return (
    <PrimaryButton
      className={variantClassName[buttonVariant]}
      onClick={handlePurchase}
    >
      Buy Now
    </PrimaryButton>
  );
};

export default BuyNowButton;
```

This component abstracts away the open feature code so that the product page component remains clean. At a high level, this button uses the `useStringFlagValue` method or hook from the Open Feature React SDK. At evaluation time, this hook will run and consult the open feature provider, which will contain the user's assignments at that moment.

From this, the evaluated variant will be used to conditionally render the black button for the control variant and the red button for the treatment variant.

You will also see that the `useStringFlagValue` hook has the `CONTROL_VARIANT` variable passed in. This is the default value. In the event that the open feature provider fails to contact the assignment service for the user's current assignments, they will be rendered the `CONTROL_VARIANT` as a fallback.

### Creating the Experiment

Now that we have instrumented the server to get events and we have set up the open feature provider to conditionally render variants on the front end, we can now create the experiment in the experimentation portal.

**1.** Open the portal and visit the experiments page.

**2.** Here we can create our experiment. Input the following values.

| Field            | Value                                                                                                                                                                                                                                                                                                                                                                                                                                                            |
| ---------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Name             | The Homer Purchase Conversion Rate Test                                                                                                                                                                                                                                                                                                                                                                                                                          |
| Feature Flag Key | the-homer-product-page-button.                                                                                                                                                                                                                                                                                                                                                                                                                                   |
| Hypothesis       | Based on prior knowledge, we believe that the current black “Buy Now” button may not be visually prominent enough for users who are ready to purchase. We think that changing the “Buy Now” button from black to red on the Power Motors Homer product page for prospective customers will achieve a higher purchase conversion rate. We will know this is true when we see a 2% improvement in purchase conversion rate compared with the black-button control. |
| Description      | This test evaluates whether or not we should change the Buy Now button to red on the homer product page.                                                                                                                                                                                                                                                                                                                                                         |

:::caution
It is vital that the feature flag key matches the feature flag declared in the open feature provider by the `BuyNowButton.tsx` file. Otherwise, variants will not be evaluated and the experiment will be compromised.
:::

**3.** Press next, and then you will be greeted with the metrics page. Here we want to define the purchase conversion rate as a binary success metric. To do this, input the following values:

| Field                           | Value              |
| ------------------------------- | ------------------ |
| Metric Key                      | purchase_conv_rate |
| Metric Type                     | Success            |
| Metric Direction                | Increase           |
| Minimum Detectable Effect (MDE) | 2                  |

Note that we have put the minimum detectable effect at 2, as we believe that changing it will increase the conversion rate by 2% as stated by our hypothesis.

**4.** Press next, and here we can define our two variants, the control and treatment. The control, our black button should be defined as the following:

| Field        | Value        |
| ------------ | ------------ |
| Variant Name | Black Button |
| Variant Key  | button-black |

Our treatment variant, the red button, should be defined as the following:

| Field        | Value      |
| ------------ | ---------- |
| Variant Name | Red Button |
| Variant Key  | button-red |

:::caution
Much like the Feature Flag key, it is vital that the variant keys match those declared in the BuyNowButton.tsx file. Otherwise, evaluations will fail, and users will just end up seeing the control variant.
:::

**5.** Congratulations, this is where the A/B test begins. This will ensure that all users see the control variant for the next seven days, starting at midnight UTC. This will allow us to establish the baseline conversion rate Which will be used as an input to our sample size calculator to inform you of the number of users that need to be in your experiment in order to accurately provide you with a recommendation at the end of the A/B test.

There is nothing more you need to do. The event metric and experiment are ready to go. Prism will collect its data in the background, and then, once these seven days are complete, you can return to the experiment portal and input other information so the A/B test can begin. This will include the percentage of the population you want to assign to the experiment and the date range of the experiment.

### Manually overriding the system to test assignment behaviour

For demonstration purposes, you likely want to check that the experimentation logic is working and users are assigned correctly. To do this, use the commands below to manually manipulate the start and end times of the experiment to move the experiment through phases.

Run the following command to start the a/a test and confirm that the user consistently sees the control variant as part of the a/a test.

```bash
docker compose exec -T postgres psql -U prism_user -d prism <<'EOF'
UPDATE prism.experiments
SET aa_start_time = now(),
    aa_end_time = now() + interval '7 day'
WHERE feature_flag_id = 'the-homer-product-page-button';
EOF
```

Then run this command to end the AA test.

```bash
docker compose exec -T postgres psql -U prism_user -d prism <<'EOF'
UPDATE prism.experiments
SET
    aa_start_time = now() - interval '7 day',
    aa_end_time = now()
WHERE feature_flag_id = 'the-homer-product-page-button';
EOF
```

Now visit the portal and you can configure the A/B test as you like. Once the form has been submitted, you can run the following command to start the experiment immediately by passing the need to wait till midnight.

```bash
docker compose exec -T postgres psql -U prism_user -d prism <<'EOF'
UPDATE prism.experiments
SET start_time = now(),
    end_time = now() + interval '7 day'
WHERE feature_flag_id = 'the-homer-product-page-button';
EOF
```

Before the correct variant will show, you will need to clear the Redis cache with the following command Which will invoke the cache invalidation logic of the assignment service

```bash
docker compose exec -T kafka bash -c 'echo "{\"action\":\"REMOVE\",\"data\":{\"experiment_key\":\"the-homer-product-page-button\",\"buckets\":[]}}" | kafka-console-producer --broker-list kafka:29092 --topic assignment-cache-invalidations'
```

If you now visit the userContext.ts file and change the `userID` value, you will see assignments change in real time. Note you might have to change the `userID` a couple of times to see the red variant.
