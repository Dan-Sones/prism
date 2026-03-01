package org.prism.eventsservice.service;

import java.util.concurrent.CompletableFuture;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.prism.eventsservice.exception.EventIngestionException;
import org.prism.eventsservice.model.EventRequest;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;
import tools.jackson.databind.ObjectMapper;

@Service
@Slf4j
@RequiredArgsConstructor
public class EventPublisher {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;

    @Value("${kafka.topic}")
    private String eventTopic;

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
                    // TODO: we can't throw an exception here as we're in a seperate thread,
                    // this should be the point that we append to a DLQ for retries.
                    // the client (sender) should get a 200 response as they have done their job.
                }
            });
        } catch (Exception e) {
            log.error("Failed to serialize event: {}", e.getMessage());
            throw new EventIngestionException("Failed to serialize event: " + e.getMessage());
        }
    }
}
