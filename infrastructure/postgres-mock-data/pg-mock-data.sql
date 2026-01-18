INSERT INTO prism.experiments (name)
VALUES
    ('Experiment A'),
    ('Experiment B'),
    ('Experiment C');


INSERT INTO prism.variants (experiment_ID, name, buckets)
VALUES
    ((SELECT id FROM prism.experiments WHERE name = 'Experiment A'), 'Variant A1', ARRAY[1,2,3]),
    ((SELECT id FROM prism.experiments WHERE name = 'Experiment A'), 'Variant A2', ARRAY[4,5]),
    ((SELECT id FROM prism.experiments WHERE name = 'Experiment B'), 'Variant B1', ARRAY[10,11]),
    ((SELECT id FROM prism.experiments WHERE name = 'Experiment C'), 'Variant C1', ARRAY[100]),
    ((SELECT id FROM prism.experiments WHERE name = 'Experiment C'), 'Variant C2', ARRAY[101,102]);
