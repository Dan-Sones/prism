ALTER TABLE prism.variants DROP COLUMN IF EXISTS buckets;

ALTER TABLE prism.variants RENAME COLUMN variant_id TO variant_key;

-- variants are no longer allocated to buckets, we will hancle that in actual code
-- some of this info seems redundant, but we can use it to visualise bucket use if I have time to develop it
CREATE TABLE IF NOT EXISTS prism.bucket_allocations (
    id SERIAL PRIMARY KEY,
    experiment_id INTEGER NOT NULL REFERENCES prism.experiments(id) ON DELETE CASCADE,
    bucket_number INTEGER NOT NULL,
    UNIQUE (experiment_id, bucket_number)
);


-- we need these anyway, but we can use thse above to visualise the experiment / bucket timeline.
ALTER TABLE prism.experiments ADD COLUMN IF NOT EXISTS start_time TIMESTAMP;
ALTER TABLE prism.experiments ADD COLUMN IF NOT EXISTS end_time TIMESTAMP;

-- we set a default of 0 for these when we alter the table and then drop it straight away - just so the migration can complete
ALTER TABLE prism.variants ADD COLUMN IF NOT EXISTS lower_bound INTEGER NOT NULL DEFAULT 0;
ALTER TABLE prism.variants ADD COLUMN IF NOT EXISTS upper_bound INTEGER NOT NULL DEFAULT 0;

ALTER TABLE prism.variants ALTER COLUMN lower_bound DROP DEFAULT;
ALTER TABLE prism.variants ALTER COLUMN upper_bound DROP DEFAULT;

ALTER TABLE prism.experiments ADD COLUMN IF NOT EXISTS unique_salt UUID UNIQUE DEFAULT gen_random_uuid();
