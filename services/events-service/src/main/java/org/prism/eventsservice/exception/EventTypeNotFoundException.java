package org.prism.eventsservice.exception;

public class EventTypeNotFoundException extends IllegalStateException {
    public EventTypeNotFoundException(String message) {
        super(message);
    }

    public EventTypeNotFoundException(String message, Throwable cause) {
        super(message, cause);
    }
}
