ALTER TABLE prism.experiments
ADD COLUMN feature_flag_id VARCHAR(255);

ALTER TABLE prism.variants
ADD COLUMN variant_id VARCHAR(255);