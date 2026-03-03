package org.prism.eventsservice.model;

import org.prism.eventsservice.grpc.events_catalog.v1.DataType;

public enum OutboundEventFieldDataType {
    STRING("string"),
    INT("int"),
    FLOAT("float"),
    BOOL("boolean"),
    TIMESTAMP("timestamp");

    private final String value;

    OutboundEventFieldDataType(String value) {
        this.value = value;
    }

    @Override
    public String toString() {
        return value;
    }

    public static OutboundEventFieldDataType fromDataType(DataType dataType) {
        return switch (dataType) {
            case DATA_TYPE_UNSPECIFIED -> null;
            case DATA_TYPE_STRING -> STRING;
            case DATA_TYPE_INT -> INT;
            case DATA_TYPE_FLOAT -> FLOAT;
            case DATA_TYPE_BOOL -> BOOL;
            case DATA_TYPE_TIMESTAMP -> TIMESTAMP;
            case UNRECOGNIZED -> throw new IllegalArgumentException("Unrecognized data type: " + dataType);
        };
    }
}
