package org.prism.eventsservice.controller;

import lombok.AllArgsConstructor;
import org.prism.eventsservice.model.EventRequest;
import org.prism.eventsservice.service.EventService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
@AllArgsConstructor
public class EventController {
    private EventService eventService;

    @PostMapping("/event")
    public ResponseEntity<Void> publishEvent(@RequestBody EventRequest eventRequest) {
        eventService.ingestEvent(eventRequest);
        return ResponseEntity.noContent().build();
    }
}
