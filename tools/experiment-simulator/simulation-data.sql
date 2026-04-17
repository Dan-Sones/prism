WITH experiment_1 AS (
    INSERT INTO prism.experiments (name, feature_flag_id, start_time, end_time)
        VALUES ('CTA Button Color Test', 'button_color_v1', NOW(), NOW() + INTERVAL '1 WEEK')
        RETURNING id),
     variant_1 AS (
         INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound, variant_type)
             VALUES ((SELECT id FROM experiment_1), 'Button - Blue', 'button_blue', 0, 100, 'control')
             RETURNING id),
     variant_2 AS (
         INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound, variant_type)
             VALUES ((SELECT id FROM experiment_1), 'Button - Green', 'button_green', 0, 0, 'treatment')
             RETURNING id),
     bucket_allocations AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             SELECT (SELECT id FROM experiment_1), g
             FROM generate_series(0, 9999) AS g),
     purchase AS (
         INSERT INTO prism.event_types (name, version, description, event_key)
             VALUES ('purchase', 1, 'Fires when a user completes a purchase', 'purchase')
             RETURNING id)
INSERT
INTO prism.event_fields (event_type_id, name, field_key, data_type)
VALUES ((SELECT id FROM purchase), 'Order Total', 'order_total', 'float');

INSERT INTO prism.metrics (name,
                           metric_key, description, metric_type, analysis_unit, is_binary)
VALUES ('Purchase Conversion Rate',
        'purchase_conv_rate',
        'The percentage of exposed users who made at
least one purchase.',
        'ratio',
        'user', true);

DO
$$
    DECLARE
        v_metric_id         UUID;
        v_purchase_event_id UUID;
        v_exposure_event_id UUID;
    BEGIN
        SELECT id
        INTO v_metric_id
        FROM prism.metrics
        WHERE metric_key = 'purchase_conv_rate';
        SELECT id
        INTO v_purchase_event_id
        FROM prism.event_types
        WHERE event_key = 'purchase';
        SELECT id
        INTO v_exposure_event_id
        FROM prism.event_types
        WHERE event_key = 'experiment_exposure';

        INSERT INTO prism.metric_components (metric_id,
                                             role, event_type_id, agg_operation, system_column_name)
        VALUES (v_metric_id, 'numerator',
                v_purchase_event_id, 'COUNT_DISTINCT', 'user_id');

        INSERT INTO prism.metric_components (metric_id,
                                             role, event_type_id, agg_operation, system_column_name)
        VALUES (v_metric_id, 'denominator',
                v_exposure_event_id, 'COUNT_DISTINCT', 'user_id');
    END
$$;
