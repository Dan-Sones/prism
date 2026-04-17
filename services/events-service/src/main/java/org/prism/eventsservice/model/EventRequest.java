package org.prism.eventsservice.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
@AllArgsConstructor
public class EventRequest {
    @JsonProperty("event_key")
    private String eventKey;

    @JsonProperty("user_details")
    private UserDetails userDetails;

    @JsonProperty("experiment_details")
    private ExperimentDetails experimentDetails;

    @JsonProperty("sent_at")
    private Instant sentAt;

    private Map<String, Object> properties;
}
