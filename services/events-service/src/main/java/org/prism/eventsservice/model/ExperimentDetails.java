package org.prism.eventsservice.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

@Data
@Builder
@AllArgsConstructor
public class ExperimentDetails {
    @JsonProperty("experiment_key")
    private String experiment_key;

    @JsonProperty("variant_key")
    private String variant_key;
}
