package org.prism.eventsservice.config;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import jakarta.annotation.PreDestroy;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class GrpcConfig {

    @Value("${grpc.server.host:localhost}")
    private String host;

    @Value("${grpc.server.port:50052}")
    private int port;

    private ManagedChannel channel;

    @Bean
    public ManagedChannel managedChannel() {
        channel = ManagedChannelBuilder.forAddress(host, port).usePlaintext().build();
        return channel;
    }

    @PreDestroy
    public void shutdown() {
        if (channel != null && !channel.isShutdown()) {
            channel.shutdown();
        }
    }
}
