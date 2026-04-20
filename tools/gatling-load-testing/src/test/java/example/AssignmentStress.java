package example;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.http;
import static io.gatling.javaapi.http.internal.HttpCheckBuilders.status;

import io.gatling.javaapi.core.ScenarioBuilder;
import io.gatling.javaapi.core.Simulation;
import io.gatling.javaapi.http.HttpProtocolBuilder;
import io.gatling.javaapi.http.HttpRequestActionBuilder;
import java.time.Duration;
import java.util.UUID;

public class AssignmentStress extends Simulation {

    HttpProtocolBuilder httpProtocol = http.baseUrl("http://35.246.99.1")
            .acceptHeader("application/json")
            .userAgentHeader(
                    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36");

    public static final HttpRequestActionBuilder getAssignments = http("Assignments")
            .get(session -> "/api/assignments/" + UUID.randomUUID())
            .check(status().is(200));
    ScenarioBuilder scen = scenario("Assignment Stress Test").exec(getAssignments);

    {
        setUp(scen.injectOpen(incrementUsersPerSec(200)
                        .times(5)
                        .eachLevelLasting(Duration.ofMinutes(2))
                        .startingFrom(200))
                .protocols(httpProtocol));
    }
}
