package org.prism.powellmotorsapispringboot.controller;

import java.time.Instant;
import java.util.Map;
import lombok.RequiredArgsConstructor;
import org.prism.powellmotorsapispringboot.client.EventsClient;
import org.prism.powellmotorsapispringboot.model.EventRequest;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api")
@RequiredArgsConstructor
public class EventController {

    private static final String DEMO_USER_ID = "demo-user-001";

    private final EventsClient eventsClient;

    @PostMapping("/view-product")
    public ResponseEntity<Void> viewProduct(@RequestParam("userId") String userId) {
        EventRequest event = EventRequest.builder()
                .eventKey("experiment_exposure")
                .experimentKey("the-homer-experiment")
                .userDetails(Map.of("id", userId))
                .sentAt(Instant.now())
                .properties(Map.of())
                .build();
        eventsClient.sendEvent(event);
        return ResponseEntity.noContent().build();
    }

    @PostMapping("/purchase")
    public ResponseEntity<Void> purchase(@RequestParam("userId") String userId) {
        EventRequest event = EventRequest.builder()
                .eventKey("purchase")
                .userDetails(Map.of("id", userId))
                .sentAt(Instant.now())
                .properties(Map.of(
                        "product_id", "the-homer",
                        "amount_gbp", 10000,
                        "order_id", "ord-123"))
                .build();
        eventsClient.sendEvent(event);
        return ResponseEntity.noContent().build();
    }
}
