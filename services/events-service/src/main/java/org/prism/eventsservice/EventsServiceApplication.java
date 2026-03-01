package org.prism.eventsservice;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.cache.annotation.EnableCaching;

@SpringBootApplication
@EnableCaching
public class EventsServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(EventsServiceApplication.class, args);
    }
}
