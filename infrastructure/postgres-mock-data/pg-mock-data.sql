INSERT INTO prism.experiments (name, feature_flag_id) VALUES
('CTA Button Color Test', 'button_color_v1');

INSERT INTO prism.variants (experiment_id, name, variant_id, buckets) VALUES
(2, 'Button - Blue', 'button_blue', ARRAY[0, 1, 2, 3, 3930]),
(2, 'Button - Green', 'button_green', ARRAY[5, 6, 7, 8, 9]);


