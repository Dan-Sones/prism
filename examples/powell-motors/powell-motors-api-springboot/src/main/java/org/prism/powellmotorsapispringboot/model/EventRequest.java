package org.prism.powellmotorsapispringboot.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import java.time.Instant;
import java.util.Map;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
@AllArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class EventRequest {
    @JsonProperty("event_key")
    private String eventKey;

    @JsonProperty("user_details")
    private Map<String, Object> userDetails;

    @JsonProperty("sent_at")
    private Instant sentAt;

    @JsonProperty("experiment_key")
    private String experimentKey;

    private Map<String, Object> properties;
}
