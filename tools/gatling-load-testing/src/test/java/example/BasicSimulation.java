package example;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.*;

import io.gatling.javaapi.core.*;
import io.gatling.javaapi.http.*;

public class BasicSimulation extends Simulation {

    // Load VU count from system properties
    // Reference: https://docs.gatling.io/guides/passing-parameters/
    private static final int virtualUsers = Integer.getInteger("virtualUsers", 10000);

    // Define HTTP configuration
    // Reference: https://docs.gatling.io/reference/script/protocols/http/protocol/
    private static final HttpProtocolBuilder httpProtocol = http.baseUrl("http://localhost:5678")
            .acceptHeader("application/json")
            .userAgentHeader(
                    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36");

    // Define scenario
    // Reference: https://docs.gatling.io/reference/script/core/scenario/
    private static final ScenarioBuilder scenario =
            scenario("Scenario").exec(http("Session").post("/event").check(status().is(200)));

    // Define injection profile and execute the test
    // Reference: https://docs.gatling.io/reference/script/core/injection/
    // Reference: https://docs.gatling.io/reference/script/core/assertions/
    {
        setUp(scenario.injectOpen(atOnceUsers(virtualUsers)))
                .assertions(
                        global().failedRequests().count().lt(1L),
                        global().responseTime().max().lt(800))
                .protocols(httpProtocol);
    }
}
