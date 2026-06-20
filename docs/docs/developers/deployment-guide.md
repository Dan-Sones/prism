---
sidebar_position: 1
---

# Deploying Prism

This documentation will help you get Prism deployed on your organisations infrastructure.

## Security Model

Before we can talk about deployment, it is first extremely important to understand the security model of the Prism platform. Prism operates under the perimeter defence model, meaning that once something is inside Prism, it is trusted and treated as truth.

## Authentication

Usually, when you think of the perimeter defense model, this would mean internet-facing services, such as REST APIs, are secured with JWT-based authentication, etc. However, as a self-hostable platform like Prism, Realistically, implementing a perimeter auth system in the assignment service and events service that fits with all user organizations' auth models is simply unfeasible. In theory, we could offer several interfaces to allow you to implement this logic yourself; however, arguably it is more efficient to move this responsibility back to you, the deploying organization.

There are two publicly facing endpoints in PRISM:

1. Within the event service, which allows for the ingestion of events
2. In the assignment service, which, given a user ID, will fetch the user's assignment at that time

These represent two fundamentally different attack vectors.

### Assignment Service Weaknesses

The assignment endpoint without authentication is vulnerable to enumeration attacks, which could result in information disclosure. In theory, an attacker could generate any number of UUIDs and then query the service to get users' assignments. However, arguably this is a non fret. Inherently, the bucketing model of prism means that, given any input, the hash function will hash that string into a bucket, and therefore assignments. This means that technically a brute force attack will actually give the attacker misleading results.

In theory, the only point at which this endpoint would offer true information leakage is if the deploying organization uses a non-UUID-based ID system for users. Or if a list of known users' UIDs were compromised

### Events Service

The events service is perhaps the most insecure service in the prism deployment. The event service takes any event ingested by a post as truth. This means the user ID contained in the request and the events attributes will all influence the metrics at the end of the experiment. This can be potentially exploited by bad actors in that if they want to influence the result of an experiment, they might choose to send fraudulent events with properties designed to influence the experiment's results. Alternatively, a sheer volume of requests for one variant or a lack of requests against a variant can be used to also influence results.

### Perimeter Defense Remediation

In order to prevent each of these attacks, neither the assignment service nor the events service should have these public-facing endpoints directly exposed to the internet. This is sort of a given anyway, as the only supported event ingestion mechanism at the moment is server-side. Consider the scenario where a user makes a purchase and you want to track purchase events. The purchase event itself should not be ingested from the client side where the user makes the purchase. The event being published should be a side effect of the server-side operations for the purchase. This means that the request has already been through your organization's authentication process, meaning that Prism can trust the user ID associated with the event.

## Deployment Model

:::danger
As discussed above, to enable the perimeter defense model, you must deploy Prism without any external internet access. Assignment requests should be forwarded to the assignment service within Prism, and events should be sent from your servers to the events service. Client devices should not interact with Prism in any way.
:::

### Docker

Each of Prism's services is deployable via Docker. Available on GitHub, the Container Registry is an image for each of the services that has been developed. Alternatively, if you wish to build the images yourselves, you can use the Makefiles located in the project to build them.

Once you have images, the recommended model of deployment is via Docker Compose and Docker networking. This allows for seamless communication between the services and external dependencies such as ClickHouse, Redis, and Postgres. For a visual representation of Prism's architecture, see below.

![Architecture Diagram](/img/arch.png)

The Docker Compose provided within the repository is pre-configured with all of the connectivity between services, And an accompanying end.example, ready to be copied, allows you to tune Prism to your needs.

See [Experiment Walkthrough](experiment-walkthrough.md) for more on how to do this.
