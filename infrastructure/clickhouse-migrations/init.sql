CREATE TABLE events (
    experiment_key    String,
    variant_key       String,
    event_key         String,
    user_id           String,
    sent_at           DateTime,
    received_at       DateTime,
    string_properties Map(String, String),
    int_properties    Map(String, Int64),
    float_properties  Map(String, Float64)
) ENGINE = MergeTree()
ORDER BY (experiment_key, variant_key, event_key, received_at);
