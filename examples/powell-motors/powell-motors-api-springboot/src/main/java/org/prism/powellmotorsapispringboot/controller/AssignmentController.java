package org.prism.powellmotorsapispringboot.controller;

import java.util.Map;
import lombok.RequiredArgsConstructor;
import org.prism.powellmotorsapispringboot.client.AssignmentClient;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api")
@RequiredArgsConstructor
public class AssignmentController {

    private final AssignmentClient assignmentClient;

    @GetMapping("/assignments/{userId}")
    public Map<String, String> getAssignments(@PathVariable String userId) {
        return assignmentClient.getAssignmentsForUser(userId);
    }
}
