WITH experiment_1 AS (
    INSERT INTO prism.experiments (name, feature_flag_id, start_time, end_time)
        VALUES ('CTA Button Color Test', 'button_color_v1', NOW(), NOW() + INTERVAL '1 WEEK')
        RETURNING id),
     variant_1 AS (
         INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
             VALUES ((SELECT id FROM experiment_1), 'Button - Blue', 'button_blue', 0, 49)
             RETURNING id),
     variant_2 AS (
         INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
             VALUES ((SELECT id FROM experiment_1), 'Button - Green', 'button_green', 50, 99)
             RETURNING id),
     bucket_allocation_1 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 1)),
     bucket_allocation_2 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 2)),
     bucket_allocation_3 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 3)),
     bucket_allocation_4 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 4)),
     bucket_allocation_5 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 5)),
     bucket_allocation_6 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 6)),
     bucket_allocation_7 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 7)),
     bucket_allocation_8 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 8)),
     bucket_allocation_9 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 9)),
     bucket_allocation_10 AS (
         INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
             VALUES ((SELECT id FROM experiment_1), 10)),
     purchase AS (
         INSERT INTO prism.event_types (name, version, description, event_key)
             VALUES ('purchase', 1, 'Fires when a user completes a purchase', 'purchase')
             RETURNING id)
INSERT
INTO prism.event_fields (event_type_id, name, field_key, data_type)
VALUES ((SELECT id FROM purchase), 'Order Total', 'order_total', 'float')


