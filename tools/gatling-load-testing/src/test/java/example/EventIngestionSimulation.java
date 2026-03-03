package example;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.*;

import io.gatling.javaapi.core.*;
import io.gatling.javaapi.http.*;

public class EventIngestionSimulation extends Simulation {

    private static final int eventCount = Integer.getInteger("eventCount", 10000);
    private static final int durationSeconds = Integer.getInteger("durationSeconds", 60);

    private static final HttpProtocolBuilder httpProtocol =
            http.baseUrl(System.getProperty("baseUrl", "http://localhost:5678")).acceptHeader("application/json");

    // This isn't really a load test, just an easy way of writing loads of requests over a period so that the
    // microbatches are written to clickhouse
    private static final ScenarioBuilder scenario = scenario("Clickhouse Write test")
            .exec(http("Event Post")
                    .post("/event")
                    .header("Content-Type", "application/json")
                    .body(StringBody("{" + "\"event_key\": \"order_shipped\","
                            + "\"user_details\": {\"id\": \"21\"},"
                            + "\"properties\": {"
                            + "\"final_order_total\": 100.00,"
                            + "\"order_total_without_discounts\": 125.99,"
                            + "\"postage_total\": 4.99"
                            + "}"
                            + "}"))
                    .check(status().is(204)));

    {
        setUp(scenario.injectOpen(constantUsersPerSec((double) eventCount / durationSeconds)
                        .during(durationSeconds)))
                .protocols(httpProtocol);
    }
}
