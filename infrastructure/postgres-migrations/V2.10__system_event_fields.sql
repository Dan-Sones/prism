ALTER TABLE prism.metric_components
    ADD COLUMN system_column_name VARCHAR(50);

ALTER TABLE prism.metric_components
    ALTER COLUMN agg_field_id DROP NOT NULL;




-- Allow metrics components to target an event_key, a system field, or JUST an event if using count
ALTER TABLE prism.metric_components
    ADD CONSTRAINT enforce_single_aggregation_target
        CHECK (
            (agg_field_id IS NOT NULL AND system_column_name IS NULL)
                OR
            (agg_field_id IS NULL AND system_column_name IS NOT NULL)
                OR
            (agg_field_id IS NULL AND system_column_name IS NULL AND agg_operation = 'COUNT')
            );