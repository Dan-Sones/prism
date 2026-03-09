WITH test_1 AS (
    INSERT INTO prism.experiments (name, feature_flag_id)
    VALUES ('Test 1', 'feature_flag_1')
    RETURNING id
),
variants_1 AS (
    INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
    VALUES
        ((SELECT id FROM test_1), 'Flag 1 - Variant 1', 'flag_1_variant_1', 0,  49),
        ((SELECT id FROM test_1), 'Flag 1 - Variant 2', 'flag_1_variant_2', 50, 99)
),
buckets_1 AS (
    INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
    SELECT id, generate_series(0, 999) FROM test_1
),

test_2 AS (
    INSERT INTO prism.experiments (name, feature_flag_id)
    VALUES ('Test 2', 'feature_flag_2')
    RETURNING id
),
variants_2 AS (
    INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
    VALUES
        ((SELECT id FROM test_2), 'Flag 2 - Variant 1', 'flag_2_variant_1', 0,  49),
        ((SELECT id FROM test_2), 'Flag 2 - Variant 2', 'flag_2_variant_2', 50, 99)
),
buckets_2 AS (
    INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
    SELECT id, generate_series(1000, 1999) FROM test_2
),

test_3 AS (
    INSERT INTO prism.experiments (name, feature_flag_id)
    VALUES ('Test 3', 'feature_flag_3')
    RETURNING id
),
variants_3 AS (
    INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
    VALUES
        ((SELECT id FROM test_3), 'Flag 3 - Variant 1', 'flag_3_variant_1', 0,  49),
        ((SELECT id FROM test_3), 'Flag 3 - Variant 2', 'flag_3_variant_2', 50, 99)
),
buckets_3 AS (
    INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
    SELECT id, generate_series(2000, 2999) FROM test_3
),

test_4 AS (
    INSERT INTO prism.experiments (name, feature_flag_id)
    VALUES ('Test 4', 'feature_flag_4')
    RETURNING id
),
variants_4 AS (
    INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
    VALUES
        ((SELECT id FROM test_4), 'Flag 4 - Variant 1', 'flag_4_variant_1', 0,  49),
        ((SELECT id FROM test_4), 'Flag 4 - Variant 2', 'flag_4_variant_2', 50, 99)
),
buckets_4 AS (
    INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
    SELECT id, generate_series(3000, 3999) FROM test_4
),

test_5 AS (
    INSERT INTO prism.experiments (name, feature_flag_id)
    VALUES ('Test 5', 'feature_flag_5')
    RETURNING id
),
variants_5 AS (
    INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
    VALUES
        ((SELECT id FROM test_5), 'Flag 5 - Variant 1', 'flag_5_variant_1', 0,  49),
        ((SELECT id FROM test_5), 'Flag 5 - Variant 2', 'flag_5_variant_2', 50, 99)
),
buckets_5 AS (
    INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
    SELECT id, generate_series(4000, 4999) FROM test_5
)

SELECT 'done' AS result;