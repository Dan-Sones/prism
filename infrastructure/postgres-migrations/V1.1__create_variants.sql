CREATE TABLE IF NOT EXISTS prism.variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    experiment_id UUID NOT NULL REFERENCES prism.experiments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    buckets INT[] NOT NULL
);