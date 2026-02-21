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

WITH page_view AS (
    INSERT INTO prism.event_types (name, version, description)
    VALUES ('page_view', 1, 'Fired when a user views a page')
    RETURNING id
),
button_click AS (
    INSERT INTO prism.event_types (name, version, description)
    VALUES ('button_click', 1, 'Fired when a user clicks a button')
    RETURNING id
)
INSERT INTO prism.event_fields (event_type_id, name, field_key, data_type)
VALUES ((SELECT id FROM page_view), 'Page URL', 'page_url', 'string'),
       ((SELECT id FROM page_view), 'Time on Page', 'time_on_page', 'float'),
       ((SELECT id FROM button_click), 'Button ID', 'button_id', 'string'),
       ((SELECT id FROM button_click), 'Click Timestamp', 'click_timestamp', 'timestamp');
