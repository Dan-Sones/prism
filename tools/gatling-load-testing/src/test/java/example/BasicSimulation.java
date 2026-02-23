package example;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.*;

import io.gatling.javaapi.core.*;
import io.gatling.javaapi.http.*;

public class BasicSimulation extends Simulation {
    // The test assumes that the application is running and the following event schemas are available in the database:

    private static final int virtualUsers = Integer.getInteger("virtualUsers", 10000);

    private static final HttpProtocolBuilder httpProtocol = http.baseUrl("http://localhost:5678")
            .acceptHeader("application/json")
            .userAgentHeader(
                    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36");

    private static final ScenarioBuilder scenario = scenario("Scenario")
            .exec(http("Event Post")
                    .post("/event")
                    .header("Content-Type", "application/json")
                    .body(StringBody("{\n" + "  \"event_key\": \"order_shipped\",\n"
                            + "  \"user_details\": {\n"
                            + "    \"id\": \"123\"\n"
                            + "  },\n"
                            + "  \"properties\": {\n"
                            + "    \"final_order_total\": 100.00,\n"
                            + "    \"order_total_without_discounts\": 125.99,\n"
                            + "    \"postage_total\": 4.99\n"
                            + "  }\n"
                            + "}"))
                    .check(status().is(200)));

    {
        setUp(scenario.injectOpen(
                        rampUsers(virtualUsers).during(30), // Ramp up
                        constantUsersPerSec(virtualUsers / 30).during(60), // Hold
                        rampUsersPerSec(virtualUsers / 30).to(0).during(30) // Ramp down
                        ))
                .assertions(
                        global().failedRequests().count().lt(1L),
                        global().responseTime().max().lt(800))
                .protocols(httpProtocol);
    }
}
