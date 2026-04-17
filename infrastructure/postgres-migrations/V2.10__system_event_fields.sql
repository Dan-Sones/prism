ALTER TABLE prism.metric_components
    ADD COLUMN system_column_name VARCHAR(50);

ALTER TABLE prism.metric_components
    ALTER COLUMN agg_field_id DROP NOT NULL;


-- there now only needs to be an aggregation on a column OR a system field like user
ALTER TABLE prism.metric_components
    ADD CONSTRAINT enforce_single_aggregation_target
        CHECK (
            (agg_field_id IS NOT NULL AND system_column_name IS NULL)
                OR
            (agg_field_id IS NULL AND system_column_name IS NOT NULL)
            );

