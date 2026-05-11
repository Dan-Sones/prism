WITH purchase AS (
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
