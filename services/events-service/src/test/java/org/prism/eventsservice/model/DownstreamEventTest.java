package org.prism.eventsservice.model;

import java.time.Instant;
import java.util.Map;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.prism.eventsservice.grpc.events_catalog.v1.DataType;
import org.prism.eventsservice.grpc.events_catalog.v1.EventField;
import org.prism.eventsservice.grpc.events_catalog.v1.EventType;

public class DownstreamEventTest {

    @Test
    void testDownstreamEventCreation() {

        String eventKey = "user_signup";

        EventType eventType = EventType.newBuilder()
                .setId("event1")
                .setEventKey(eventKey)
                .setDescription("User signup event")
                .addFields(EventField.newBuilder()
                        .setId("field1")
                        .setFieldKey("username")
                        .setName("username")
                        .setDataType(DataType.DATA_TYPE_STRING)
                        .build())
                .addFields(EventField.newBuilder()
                        .setId("field2")
                        .setFieldKey("age")
                        .setName("age")
                        .setDataType(DataType.DATA_TYPE_INT)
                        .build())
                .addFields(EventField.newBuilder()
                        .setId("field3")
                        .setFieldKey("cost")
                        .setName("cost")
                        .setDataType(DataType.DATA_TYPE_FLOAT)
                        .build())
                .addFields(EventField.newBuilder()
                        .setId("field4")
                        .setFieldKey("signup_time")
                        .setName("signup_time")
                        .setDataType(DataType.DATA_TYPE_TIMESTAMP)
                        .build())
                .addFields(EventField.newBuilder()
                        .setId("field5")
                        .setFieldKey("is_premium")
                        .setName("is_premium")
                        .setDataType(DataType.DATA_TYPE_BOOL)
                        .build())
                .build();

        UserDetails userDetails = new UserDetails("user123");
        Map<String, Object> properties = Map.of(
                "username",
                "john_doe",
                "age",
                30,
                "cost",
                19.99,
                "signup_time",
                System.currentTimeMillis(),
                "is_premium",
                true);
        Instant sentAt = Instant.now();
        EventRequest event = EventRequest.builder()
                .eventKey(eventKey)
                .userDetails(userDetails)
                .sentAt(sentAt)
                .properties(properties)
                .build();

        DownstreamEvent downstreamEvent = new DownstreamEvent(eventType, event);

        Assertions.assertEquals(eventType.getId(), downstreamEvent.getId());
        Assertions.assertEquals(eventType.getEventKey(), downstreamEvent.getEventKey());

        event.getProperties().forEach((key, value) -> {
            Assertions.assertEquals(
                    value, downstreamEvent.getProperties().get(key).getValue());
        });

        Assertions.assertEquals(
                OutboundEventFieldDataType.STRING,
                downstreamEvent.getProperties().get("username").getDataType());
        Assertions.assertEquals(
                OutboundEventFieldDataType.INT,
                downstreamEvent.getProperties().get("age").getDataType());
        Assertions.assertEquals(
                OutboundEventFieldDataType.FLOAT,
                downstreamEvent.getProperties().get("cost").getDataType());
        Assertions.assertEquals(
                OutboundEventFieldDataType.TIMESTAMP,
                downstreamEvent.getProperties().get("signup_time").getDataType());
        Assertions.assertEquals(
                OutboundEventFieldDataType.BOOL,
                downstreamEvent.getProperties().get("is_premium").getDataType());

        Assertions.assertEquals(
                userDetails.getId(), downstreamEvent.getUserDetails().getId());

        Assertions.assertEquals(sentAt, downstreamEvent.getSentAt());
        Assertions.assertNotNull(downstreamEvent.getReceivedAt());
    }
}
