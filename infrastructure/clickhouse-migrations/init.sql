CREATE TABLE events (
    event_key         String,
    user_id           String,
    sent_at           DateTime,
    received_at       DateTime,
    string_properties Map(String, String),
    int_properties    Map(String, Int64),
    float_properties  Map(String, Float64)
) ENGINE = MergeTree()
ORDER BY (received_at);