CREATE TABLE status(
    id        INTEGER PRIMARY KEY,
    last_seen timestamp NOT NULL DEFAULT NOW()
)
