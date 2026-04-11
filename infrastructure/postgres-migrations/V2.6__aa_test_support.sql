ALTER TABLE prism.experiments ADD aa_start_time TIMESTAMP;
ALTER TABLE prism.experiments ADD aa_end_time TIMESTAMP;

ALTER TABLE prism.experiment_metric ADD COLUMN variance FLOAT;