---
sidebar_position: 3
---

# Assignment Service

The Assignment Service is a lightweight but core component of the Prism Experimentation Platform.

It is responsible for fetching the experiment assignments when they interact with your application.

This is achieved via a GRPC call to the Admin Service, which returns the appropriate assignments based on the user's ID, and therefore bucket.

## API Specification

[View Assignment Service API Documentation →](/api/assignment-service)
