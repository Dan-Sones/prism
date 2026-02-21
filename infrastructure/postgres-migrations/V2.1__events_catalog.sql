CREATE TABLE IF NOT EXISTS prism.event_types
(
    id          UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    name        VARCHAR(255) NOT NULL CONSTRAINT unique_event_type_name UNIQUE,
    version     INTEGER      NOT NULL DEFAULT 1,-- I won't be using this for now, it will always be 1. But if we allow schema evolution it will be tracked here.
    description TEXT,
    created_at  TIMESTAMPTZ           DEFAULT NOW()
);

CREATE TYPE prism.data_type_enum AS ENUM ('string', 'int', 'float', 'boolean', 'timestamp');

CREATE TABLE IF NOT EXISTS prism.event_fields
(
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type_id UUID REFERENCES prism.event_types (id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    field_key     VARCHAR(255) NOT NULL,
    data_type     prism.data_type_enum NOT NULL,
    CONSTRAINT unique_event_type_field_key UNIQUE (event_type_id, field_key)
);
