INSERT INTO prism.experiments (name, feature_flag_id)
VALUES ('CTA Button Color Test', 'button_color_v1');

INSERT INTO prism.variants (experiment_id, name, variant_key, lower_bound, upper_bound)
VALUES (1, 'Button - Blue', 'button_blue', 0, 49),
       (1, 'Button - Green', 'button_green', 50, 99);

INSERT INTO prism.bucket_allocations (experiment_id, bucket_number)
VALUES (1, 0),
       (1, 1),
       (1, 2),
       (1, 3),
       (1, 3930),
       (1, 5),
       (1, 6),
       (1, 7),
       (1, 8),
       (1, 9);
