CREATE TYPE prism.variant_type AS ENUM ('control', 'treatment');

ALTER TABLE prism.variants ADD COLUMN variant_type prism.variant_type NOT NULL;