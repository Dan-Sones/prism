CREATE TABLE IF NOT EXISTS prism.variants (
    id SERIAL PRIMARY KEY,
    experiment_id INTEGER NOT NULL REFERENCES prism.experiments(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    buckets INT[] NOT NULL
);