ALTER TABLE prism.event_types
    ADD COLUMN event_key VARCHAR(50) NOT NULL CONSTRAINT unique_event_type_event_key UNIQUE;
