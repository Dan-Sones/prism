package org.prism.eventsservice.controller;

import org.prism.eventsservice.exception.EventIngestionException;
import org.prism.eventsservice.exception.EventTypeNotFoundException;
import org.prism.eventsservice.model.EventIngestionFailure;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

@ControllerAdvice
public class ErrorController {

    @ExceptionHandler(EventIngestionException.class)
    public ResponseEntity<EventIngestionFailure> handleEventIngestionException(EventIngestionException ex) {
        String error = "Event ingestion failed: " + ex.getMessage();
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(new EventIngestionFailure(error));
    }

    @ExceptionHandler(EventTypeNotFoundException.class)
    public ResponseEntity<EventIngestionFailure> handleEventTypeNotFoundException(EventTypeNotFoundException ex) {
        String error = "Event type not found. Make sure you have set the correct event key. " + ex.getMessage();
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(new EventIngestionFailure(error));
    }
}
