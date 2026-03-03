package org.prism.eventsservice.model;

import java.time.Instant;
import java.util.Map;
import java.util.stream.Collectors;
import lombok.Data;
import org.prism.eventsservice.grpc.events_catalog.v1.EventField;
import org.prism.eventsservice.grpc.events_catalog.v1.EventType;

@Data
public class DownstreamEvent {

    private String id;
    private String eventKey;
    private UserDetails userDetails;
    private Instant sentAt;
    private Instant receivedAt;
    private Map<String, OutboundEventField> properties;

    public DownstreamEvent(EventType eventDefinition, EventRequest eventRequest) {
        this.id = eventDefinition.getId();
        this.eventKey = eventDefinition.getEventKey();
        this.userDetails = eventRequest.getUserDetails();
        this.sentAt = eventRequest.getSentAt();
        this.receivedAt = Instant.now();

        // By streaming like this events without definitions will be ignored and NOT written to kafka
        this.properties = eventDefinition.getFieldsList().stream()
                .collect(Collectors.toMap(EventField::getFieldKey, eventField -> {
                    Object value = eventRequest.getProperties().get(eventField.getFieldKey());
                    return combinePropertiesAndDefinitionsForEvent(eventField, value);
                }));
    }

    private OutboundEventField combinePropertiesAndDefinitionsForEvent(EventField eventField, Object value) {
        OutboundEventField outboundEventField = new OutboundEventField();
        outboundEventField.setDataType(OutboundEventFieldDataType.fromDataType(eventField.getDataType()));
        outboundEventField.setValue(value);
        return outboundEventField;
    }
}
