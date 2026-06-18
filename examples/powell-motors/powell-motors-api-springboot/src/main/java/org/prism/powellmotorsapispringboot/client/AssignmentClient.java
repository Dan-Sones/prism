package org.prism.powellmotorsapispringboot.client;

import java.util.Map;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.stereotype.Component;
import org.springframework.web.client.RestClient;

@Component
@RequiredArgsConstructor
public class AssignmentClient {

    private static final ParameterizedTypeReference<Map<String, String>> RESPONSE_TYPE =
            new ParameterizedTypeReference<>() {};

    private final RestClient restClient;

    @Value("${assignment-service.url:http://localhost:8082}")
    private String baseUrl;

    public Map<String, String> getAssignmentsForUser(String userId) {
        return restClient
                .get()
                .uri(baseUrl + "/api/assignments/{userId}", userId)
                .retrieve()
                .body(RESPONSE_TYPE);
    }
}
