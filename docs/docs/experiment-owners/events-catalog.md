---
sidebar_position: 2
---

# Events Catalog

Events are actions your user performs which prism records. These events them form the basis of Metrics.

Example events may include:

- a "see more" button click
- a purchase event
- an order shipment

Events can only be recorded, stored, and used for analysis in Prism if they are correctly configured in the Events Catalog.

**Important Note:** Client-Side and Server-Side events will be **DROPPED** if they are not defined in the Events Catalog.

## Creating a New Event Type

Event Types define one of the actions provided as examples above.

The process for defining a Client-Side Event _(an event delivered to prism from the frontend)_ or a Server-Side Event _(an event delivered to prism from a backend service)_ is the same.

### Instructions

1. Visit the /events-catalog/create page. You should see the following form:

![Create Event Form](/img/create-event-type-form.png)

2. **Name**: Choose a name that is easily identifies your event specifically. For example the name _"Button Press"_ would be considered weak, instead choose something more specific that uniquely identifies your event such as _"Client-Side Product Page See More Button Press"_.

3. **Event Key**: Choosing the Event Key is incredibly important. This **UNIQUE** key identifies your event globally across prism. If this key does not exactly match the key sent in Client-Side or Server-Side events, then Prism will have no way of ingesting the event. Make sure to collaborate with developers to choose this name.

4. **Description**: The description field is not mandatory, but it's useful for other experiment owners to understand the what, when, and how of your event.

5. **Fields**: Event Fields are the key:value pairs within your event that contain properties used to form metrics. In the case of a _Purchase_ event, this might be:

   ```json
   {
     "orderTotal": 14.99,
     "postageCost": 4.99,
     "couponUsed": false
   }
   ```

   For each field you must input the following:

   | Field     | Description                                                                                                                                                                            |
   | --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
   | Name      | Display name of the field (to be used in the ui)                                                                                                                                       |
   | Field Key | The key used when sending this field in events. Again it's **Vital** this matches the exact key sent within the event itself. Collaborate with developers if required to make this so. |
   | DataType  | The data type of the field (e.g., `string`, `float`, `integer`, `boolean`). If this data type is not inputted corretly the then metrics maybe malformed and experiment results invalid |

6. Once Create Event has been pressed, the event is ready for usage across the prism architecture.


## Event Details

