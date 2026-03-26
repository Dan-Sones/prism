package org.prism.eventsservice;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;

@SpringBootApplication
@EnableCaching
@Slf4j
public class EventsServiceApplication {
    public static void main(String[] args) {
        log.info("Starting Events Service Application...");
        SpringApplication.run(EventsServiceApplication.class, args);
    }
}
