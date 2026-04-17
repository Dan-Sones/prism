package org.prism.eventsservice.model;

import java.util.List;

public record EventValidationResult(boolean isValid, List<String> missingFields) {}
