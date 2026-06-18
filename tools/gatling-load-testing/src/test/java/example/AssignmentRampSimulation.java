package example;

import static io.gatling.javaapi.core.CoreDsl.*;
import static io.gatling.javaapi.http.HttpDsl.*;

import io.gatling.javaapi.core.*;
import io.gatling.javaapi.http.*;
import java.time.Duration;
import java.util.*;

public class AssignmentRampSimulation extends Simulation {
    HttpProtocolBuilder httpProtocol = http.baseUrl("http://localhost:8082").acceptHeader("application/json");

    ScenarioBuilder scn = scenario("Ramp scenario")
            .exec(http("Get endpoint")
                    .get(session -> "/api/assignments/" + UUID.randomUUID().toString()));
    ;

    {
        setUp(scn.injectOpen(
                incrementUsersPerSec(50)
                        .times(10)
                        .eachLevelLasting(Duration.ofSeconds(90))
                        .separatedByRampsLasting(Duration.ofSeconds(10))
                        .startingFrom(200)
        ))
                .protocols(httpProtocol);
    }
}
