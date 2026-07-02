package org.prism.eventsservice.service;

import io.grpc.ManagedChannel;
import io.grpc.Status;
import io.grpc.StatusRuntimeException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.exception.EventIngestionException;
import org.prism.eventsservice.grpc.events_catalog.v1.DataType;
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
            throw new EventIngestionException(
                    "Missing required fields: " + String.join(", ", eventValidationResult.missingFields()));
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
            throw new EventIngestionException(
                    "Invalid event properties: " + String.join(", ", propertiesValidationResult.validationErrors()));
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

        if (eventRequest.getUserDetails().getId() == null
                || eventRequest.getUserDetails().getId().isEmpty()) {
            missingFields.add("userDetails.id");
        }

        if (eventRequest.getSentAt() == null) {
            missingFields.add("sentAt");
        }

        if ("experiment_exposure".equals(eventRequest.getEventKey())
                && (eventRequest.getExperimentKey() == null
                        || eventRequest.getExperimentKey().isEmpty())) {
            missingFields.add("experimentKey");
        }

        return new EventValidationResult(missingFields.isEmpty(), missingFields);
    }

    private EventPropertiesValidationResult validateEventProperties(
            Map<String, Object> eventProperties, EventType eventType) {
        ArrayList<String> validationErrors = new ArrayList<>();

        if (eventProperties == null) {
            validationErrors.add("properties must be present");
            // TODO: Surface error in portal that event was sent with no properties
            return new EventPropertiesValidationResult(false, validationErrors);
        }

        for (var schemaField : eventType.getFieldsList()) {
            String fieldKey = schemaField.getFieldKey();

            if (!eventProperties.containsKey(fieldKey)) {
                validationErrors.add("Missing property: " + fieldKey);
                // TODO: Surface error in portal that event is missing property
                continue;
            }

            Object value = eventProperties.get(fieldKey);

            if (!isCorrectType(value, schemaField.getDataType())) {
                // TODO: Surface error in portal that there is a type mismatch
                validationErrors.add("Property " + fieldKey + " expected " + schemaField.getDataType() + " but got "
                        + value.getClass().getSimpleName());
            }
        }

        return new EventPropertiesValidationResult(validationErrors.isEmpty(), validationErrors);
    }

    private boolean isCorrectType(Object value, DataType expectedType) {
        if (value == null) {
            return false;
        }
        return switch (expectedType) {
            case DATA_TYPE_STRING -> value instanceof String;
            case DATA_TYPE_INT -> value instanceof Integer;
            case DATA_TYPE_FLOAT -> value instanceof Number;
            case DATA_TYPE_BOOL -> value instanceof Boolean;
            case DATA_TYPE_TIMESTAMP -> value instanceof String;
            default -> false;
        };
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
