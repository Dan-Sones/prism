package org.prism.eventsservice.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.Map;
import lombok.Data;

@Data
public class EventRequest {
    @JsonProperty("event_key")
    private String eventKey;

    @JsonProperty("user_details")
    private UserDetails userDetails;

    private Map<String, Object> properties;
}
