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


-- This is pretty denormalised, we are duplicating almost everything. I'm not sure if this is a good idea.
-- TODO: Look into whether this table should ref the events table?
CREATE TABLE cooked_events (
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
