package org.prism.eventsservice.service;

import io.grpc.ManagedChannel;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.exception.EventIngestionException;
import org.prism.eventsservice.grpc.events_catalog.v1.EventType;
import org.prism.eventsservice.grpc.events_catalog.v1.EventsCatalogServiceGrpc;
import org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest;
import org.prism.eventsservice.model.DownstreamEvent;
import org.prism.eventsservice.model.EventPropertiesValidationResult;
import org.prism.eventsservice.model.EventRequest;
import org.prism.eventsservice.model.EventValidationResult;
import org.springframework.cache.CacheManager;
import org.springframework.stereotype.Service;

@Service
@Slf4j
public class EventService {

    private final EventsCatalogServiceGrpc.EventsCatalogServiceBlockingStub eventsCatalogStub;
    private final CacheManager cacheManager;
    private final EventPublisher eventPublisher;

    public EventService(ManagedChannel channel, CacheManager cacheManager, EventPublisher eventPublisher) {
        this.eventsCatalogStub = EventsCatalogServiceGrpc.newBlockingStub(channel);
        this.cacheManager = cacheManager;
        this.eventPublisher = eventPublisher;
    }

    public void ingestEvent(EventRequest eventToIngest) {

        var eventValidationResult = validateEvent(eventToIngest);
        if (!eventValidationResult.isValid()) {
            throw new EventIngestionException("Missing required fields: " + String.join(", ", eventValidationResult.missingFields()));
        }

        EventType eventType;
        try {
            eventType = lookupEventType(eventToIngest.getEventKey());
        } catch (Exception e) {
            log.error("Failed to lookup event type for key {}: {}", eventToIngest.getEventKey(), e.getMessage());
            // we failed internally, not the clients fault, we should send them a 200 and we deal with retrying later on
            if (e instanceof EventIngestionException) {
                // we need to report this to the user as they are sending events with misconfigure event keys!!
                throw e;
            }
            return;
        }

        var propertiesValidationResult = validateEventProperties(eventToIngest.getProperties(), eventType);
        if (!propertiesValidationResult.isValid()) {
            return;
        }

        DownstreamEvent downstreamEvent = new DownstreamEvent(eventType, eventToIngest);

        try {
            eventPublisher.publish(downstreamEvent);
        } catch (Exception e) {
            log.error("Failed to publish event: {}", e.getMessage());
            throw e;
        }
    }

    private EventValidationResult validateEvent(EventRequest eventRequest) {

        List<String> missingFields = new ArrayList<>();

        if (eventRequest.getEventKey() == null || eventRequest.getEventKey().isEmpty()) {
            missingFields.add("eventKey");
        }

        if (eventRequest.getUserDetails().getId() == null || eventRequest.getUserDetails().getId().isEmpty()) {
            missingFields.add("userDetails.id");
        }

        if (eventRequest.getSentAt() == null) {
            missingFields.add("sentAt");
        }

        if (eventRequest.getExperimentDetails().getExperiment_key() == null || eventRequest.getExperimentDetails().getExperiment_key().isEmpty()) {
            missingFields.add("experimentDetails.experimentKey");
        }

        if (eventRequest.getExperimentDetails().getVariant_key() == null || eventRequest.getExperimentDetails().getVariant_key().isEmpty()) {
            missingFields.add("experimentDetails.variantKey");
        }

        return new EventValidationResult(missingFields.isEmpty(), missingFields);
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

    public EventType lookupEventType(String eventKey) {
        var cachedEventType = cacheManager.getCache("eventTypes").get(eventKey, EventType.class);
        if (cachedEventType != null) {
            return cachedEventType;
        }

        try {
            var eventType = eventsCatalogStub
                    .getEventTypeByKey(GetEventTypeByKeyRequest.newBuilder()
                            .setEventKey(eventKey)
                            .build())
                    .getEventType();
            cacheManager.getCache("eventTypes").put(eventKey, eventType);
            return eventType;
        } catch (StatusRuntimeException e) {
            if (e.getStatus().getCode() == Status.Code.NOT_FOUND) {
                throw new EventIngestionException("Event type not found for key: " + eventKey, e);
            }
            throw e;
        }
    }
}
