CREATE TYPE prism.metric_role AS ENUM ('success', 'guardrail', 'deterioration', 'quality');
CREATE TYPE prism.metric_direction AS ENUM('increase', 'decrease');

CREATE TABLE prism.experiment_metric
(
    metric_id     UUID REFERENCES prism.metrics (id) ON UPDATE CASCADE ON DELETE CASCADE,
    experiment_id UUID REFERENCES prism.experiments (id) ON UPDATE CASCADE ON DELETE CASCADE,
    role          prism.metric_role      NOT NULL,
    direction     prism.metric_direction NOT NULL,
    mde           FLOAT,
    nim           FLOAT,
    PRIMARY KEY (metric_id, experiment_id, role),

    CONSTRAINT chk_mde_success CHECK (
        (role = 'success' AND mde IS NOT NULL) OR
        (role != 'success' AND mde IS NULL)
        ),

    CONSTRAINT chk_nim_guardrail CHECK (
        (role = 'guardrail' AND nim IS NOT NULL) OR
        (role != 'guardrail' AND nim IS NULL)
        )
);

