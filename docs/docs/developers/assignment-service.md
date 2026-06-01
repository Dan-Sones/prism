---
sidebar_position: 3
---

# Assignment Service

The Assignment Service responds to requests containing a user ID, and provides the user with their up-to-date experiment assignments. This service sits in the critical request path of any page in your site that is running an experiment.

## Interaction Model

In most cases you will never call the assignment service directly. This behavior is abstracted away behind the Prism Open Feature Provider, see [Open Feature](./open-feature.md).

You may call the directly in circumstances such as:

- You are doing server-side assignment evaluation
- You are integrating into a system that does not / cannot use OpenFeature.

## Key Ideas

### Buckets

In Prism experiments are assigned to buckets, and buckets contain users. Buckets are stateless but deterministic. The same user ID provided to the bucketing hash function will always output the same bucket. This means there is no need to store a database row stating a user is assigned to a bucket, it can just be calculated at run time.

:::danger
The number of buckets is FIXED. As a developer when first setting up Prism you can choose the number of buckets you want with the `GLOBAL_BUCKET_COUNT` environment variable. Once set it CANNOT be changed as it will reshuffle users into different buckets, spoiling any experiments.
:::

### Variant Assignment (Numberline Method)

Buckets only hash users to experiments NOT to variants. A users variant is determined by hashing their user ID with an experiment specific salt to produce a value 0-99. This is then mapped against assignment bounds as illustrated below on a number line. Where the user falls on this numberline determines the variant they see. Like buckets this hash function is deterministic.

![Numberline Visualized](/img/75-25.png)

### Caching

Despite assignments at both the experiment and variant level being stateless for _users_ buckets assignment to experiments still need to be stored in a database. Experiment variant bounds also.

Fetching assignments from the assignment service is a high volume activity. For each assignment request the service must fetch the active experiments for that bucket via gRPC from the experimentation service.

Experiment details cannot change, and experiments are always active from and to 00:00 UTC of given dates. This means per-day the cache aside pattern can be taken advantage of. Each day, as a user within x bucket makes a request and the cache entry for bucket x is empty, the request can be made to the experimentation-service, and the resultant data can be stored in the cache, preventing future requests.

## Process Diagram

![Assignment Diagram](/img/assignment-diagram.png)

## API Specification

[View Assignment Service API Documentation →](/api/assignment-service)
