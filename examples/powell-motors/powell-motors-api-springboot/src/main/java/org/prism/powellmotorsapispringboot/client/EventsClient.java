package org.prism.powellmotorsapispringboot.client;

import lombok.RequiredArgsConstructor;
import org.prism.powellmotorsapispringboot.model.EventRequest;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;

@Component
@RequiredArgsConstructor
public class EventsClient {

    private final RestClient restClient;

    @Value("${events-service.url:http://localhost:5678}")
    private String baseUrl;

    public void sendEvent(EventRequest event) {
        restClient.post().uri(baseUrl + "/event").body(event).retrieve().toBodilessEntity();
    }
}
