package org.prism.eventsservice.service;

import io.grpc.ManagedChannel;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.grpc.events_catalog.v1.EventType;
import org.prism.eventsservice.grpc.events_catalog.v1.EventsCatalogServiceGrpc;
import org.prism.eventsservice.grpc.events_catalog.v1.GetEventTypeByKeyRequest;
import org.prism.eventsservice.model.EventRequest;
import org.springframework.stereotype.Service;

@Service
@Slf4j
public class EventService {

    private final EventsCatalogServiceGrpc.EventsCatalogServiceBlockingStub eventsCatalogStub;

    public EventService(ManagedChannel channel) {
        this.eventsCatalogStub = EventsCatalogServiceGrpc.newBlockingStub(channel);
    }

    public void IngestEvent(EventRequest eventRequest) {
        log.info("Event Ingested Bruv" + " " + eventRequest.getEventKey() + " "
                + eventRequest.getUserDetails().getId());

        EventType eventType = eventsCatalogStub.getEventTypeByKey(GetEventTypeByKeyRequest.newBuilder()
                .setEventKey(eventRequest.getEventKey())
                .build());

        log.info("Event Type Retrieved: " + eventType.getName());
    }
}
