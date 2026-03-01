package org.prism.eventsservice.model;

import java.util.List;

public record EventPropertiesValidationResult(boolean isValid, List<String> missingFields) {}
