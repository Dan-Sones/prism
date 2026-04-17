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
             RETURNING id),
     exposure AS (
         INSERT INTO prism.event_types (name, version, description, event_key)
             VALUES ('experiment_exposure', 1, 'Fires when a user is exposed to an experiment variant', 'experiment_exposure')
             RETURNING id)
INSERT
INTO prism.event_fields (event_type_id, name, field_key, data_type)
VALUES ((SELECT id FROM purchase), 'Order Total', 'order_total', 'float')
