CREATE TYPE prism.metric_type AS ENUM ('simple', 'ratio');
CREATE TYPE prism.analysis_unit AS ENUM ('user_id');
CREATE TYPE prism.component_role AS ENUM ('base_event', 'numerator', 'denominator');
CREATE TYPE prism.aggregation_operation AS ENUM ('COUNT', 'SUM', 'AVG', 'MIN', 'MAX', 'COUNT_DISTINCT', 'PERCENTILE_95', 'PERCENTILE_99');



CREATE TABLE IF NOT EXISTS prism.metrics
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          VARCHAR(255)        NOT NULL
        CONSTRAINT unique_metric_name UNIQUE,
    metric_key    VARCHAR(50)         NOT NULL
        CONSTRAINT unique_metric_key UNIQUE,
    description   TEXT,
    created_at    TIMESTAMPTZ      DEFAULT NOW(),
    metric_type   prism.metric_type   NOT NULL,
    analysis_unit prism.analysis_unit NOT NULL
);

ALTER TABLE prism.event_fields
    ADD CONSTRAINT unique_event_field_with_type UNIQUE (id, event_type_id);

CREATE TABLE prism.metric_components
(
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_id          UUID                        NOT NULL REFERENCES prism.metrics (id) ON DELETE CASCADE,
    role               prism.component_role        NOT NULL,
    event_type_id      UUID                        NOT NULL REFERENCES prism.event_types (id),
    agg_operation prism.aggregation_operation NOT NULL,
    agg_field_id  UUID,

    FOREIGN KEY (agg_field_id, event_type_id) REFERENCES prism.event_fields (id, event_type_id)
);

