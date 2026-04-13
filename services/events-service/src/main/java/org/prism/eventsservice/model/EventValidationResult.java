package org.prism.eventsservice.model;

import java.util.List;
import java.util.Map;

public record EventValidationResult(boolean isValid, List<String> missingFields) {}
