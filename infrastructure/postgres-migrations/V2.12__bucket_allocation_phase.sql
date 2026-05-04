CREATE TYPE prism.experiment_phase AS ENUM ('AA', 'AB');

ALTER TABLE prism.bucket_allocations
    ADD COLUMN phase prism.experiment_phase NOT NULL DEFAULT 'AA';

ALTER TABLE prism.bucket_allocations
    DROP CONSTRAINT bucket_allocations_experiment_id_bucket_number_key;

ALTER TABLE prism.bucket_allocations
    ADD CONSTRAINT bucket_allocations_experiment_id_bucket_number_phase_key
        UNIQUE (experiment_id, bucket_number, phase);

ALTER TABLE prism.bucket_allocations
    ALTER COLUMN phase DROP DEFAULT;
