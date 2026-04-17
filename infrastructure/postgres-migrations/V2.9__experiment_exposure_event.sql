-- TODO: Add event types that cannot be deleted

INSERT INTO prism.event_types (name, event_key, version, description)
    VALUES ('Experiment Exposure', 'experiment_exposure', 1,
            'Fired when a user is exposed to an experiment.');
