package org.prism.eventsservice.model;

import lombok.Data;

@Data
public class OutboundEventField {
    private Object value;
    private OutboundEventFieldDataType dataType;
}
