package org.prism.eventsservice.service;

import io.grpc.ManagedChannel;
import java.util.ArrayList;
import java.util.Map;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.grpc.events_catalog.v1.EventType;
import org.prism.eventsservice.grpc.events_catalog.v1.EventsCatalogServiceGrpc;
import org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest;
import org.prism.eventsservice.model.EventPropertiesValidationResult;
import org.prism.eventsservice.model.EventRequest;
import org.springframework.stereotype.Service;

@Service
@Slf4j
public class EventService {

    private final EventsCatalogServiceGrpc.EventsCatalogServiceBlockingStub eventsCatalogStub;

    public EventService(ManagedChannel channel) {
        this.eventsCatalogStub = EventsCatalogServiceGrpc.newBlockingStub(channel);
    }

    public void IngestEvent(EventRequest eventToIngest) {
        var eventType = lookupEventType(eventToIngest.getEventKey());

        var validationResult = validateEventProperties(eventToIngest.getProperties(), eventType);
        if (!validationResult.isValid()) {
            return;
        }

        log.info(
                "Ingesting event of type " + eventType.getName() + " with properties " + eventToIngest.getProperties());
    }

    private EventPropertiesValidationResult validateEventProperties(
            Map<String, Object> eventProperties, EventType eventType) {
        ArrayList<String> missingFields = new ArrayList<>();

        for (var schemaField : eventType.getFieldsList()) {
            if (!eventProperties.containsKey(schemaField.getFieldKey())) {
                // TODO: Missing fields need to be raised through an observable alert in the portal, for now just log
                // and discard the event
                log.warn("Missing field " + schemaField.getFieldKey() + " in event " + eventType.getName());
                missingFields.add(schemaField.getFieldKey());
            }
        }

        return new EventPropertiesValidationResult(missingFields.isEmpty(), missingFields);
    }

    private EventType lookupEventType(String eventKey) {
        // TODO: First look in the redis cache, if not found then make a GRPC call and cache-aside
        return eventsCatalogStub.getEventTypeByKey(
                GetEventTypeByKeyRequest.newBuilder().setEventKey(eventKey).build());
    }
}
