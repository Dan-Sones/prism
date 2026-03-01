package org.prism.eventsservice.exception;

public class EventIngestionException extends RuntimeException {
    public EventIngestionException(String message) {
        super(message);
    }

    public EventIngestionException(String message, Throwable cause) {
        super(message, cause);
    }
}
