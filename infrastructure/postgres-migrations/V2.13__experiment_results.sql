CREATE TYPE prism.decision_recommendation AS ENUM (
    'DECISION_RECOMMENDATION_UNSPECIFIED',
    'DECISION_RECOMMENDATION_RECOMMEND',
    'DECISION_RECOMMENDATION_NOT_RECOMMEND',
    'DECISION_RECOMMENDATION_INCONCLUSIVE'
);

CREATE TABLE prism.experiment_results
(
    experiment_id         UUID PRIMARY KEY REFERENCES prism.experiments (id) ON DELETE CASCADE,
    recommendation        prism.decision_recommendation NOT NULL DEFAULT 'DECISION_RECOMMENDATION_UNSPECIFIED',
    recommendation_reason TEXT,
    calculated_at         TIMESTAMP
);

CREATE TABLE prism.ztest_results
(
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experiment_id             UUID              NOT NULL,
    metric_id                 UUID              NOT NULL,
    role                      prism.metric_role NOT NULL,
    absolute_difference       BOOLEAN           NOT NULL,
    ci_lower                  DOUBLE PRECISION  NOT NULL,
    ci_upper                  DOUBLE PRECISION  NOT NULL,
    p_value                   DOUBLE PRECISION  NOT NULL,
    adjusted_ci_lower         DOUBLE PRECISION,
    adjusted_ci_upper         DOUBLE PRECISION,
    adjusted_p_value          DOUBLE PRECISION,
    is_significant            BOOLEAN           NOT NULL,
    powered_effect            DOUBLE PRECISION,
    created_at                TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
--     I'm not sure if these belong here - they are related to the z test (the input) - but they are not strictly results of the z test
    control_numerator         BIGINT            NOT NULL,
    control_denominator       BIGINT            NOT NULL,
    treatment_numerator       BIGINT            NOT NULL,
    treatment_denominator     BIGINT            NOT NULL,
    practically_significant   BOOLEAN           NOT NULL,
    statistically_significant BOOLEAN           NOT NULL,
    CONSTRAINT fk_experiment_metric
        FOREIGN KEY (metric_id, experiment_id, role)
            REFERENCES prism.experiment_metric (metric_id, experiment_id, role)
            ON UPDATE CASCADE ON DELETE CASCADE
);