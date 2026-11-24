CREATE TABLE IF NOT EXISTS items (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(200) NOT NULL,
    quantity    INT NOT NULL DEFAULT 0,
    location_id INT REFERENCES locations(id),
    deleted_at  TIMESTAMPTZ
);
