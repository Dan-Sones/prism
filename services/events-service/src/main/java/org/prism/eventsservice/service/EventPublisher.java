package org.prism.eventsservice.service;

import java.util.concurrent.CompletableFuture;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.model.EventRequest;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import tools.jackson.databind.ObjectMapper;

@Service
@Slf4j
public class EventPublisher {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Value("${kafka.topic}")
    private String eventTopic;

    public EventPublisher(KafkaTemplate<String, String> kafkaTemplate, ObjectMapper objectMapper) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
    }

    public void publish(EventRequest event) {
        try {
            String json = objectMapper.writeValueAsString(event);
            CompletableFuture<SendResult<String, String>> future = kafkaTemplate.send(eventTopic, json);
            future.whenComplete((result, ex) -> {
                if (ex == null) {
                    log.info(
                            "Event published successfully with offset=[{}]",
                            result.getRecordMetadata().offset());
                } else {
                    log.error("Failed to publish event: {}", ex.getMessage());
                }
            });
        } catch (Exception e) {
            log.error("Failed to serialize event: {}", e.getMessage());
        }
    }
}
