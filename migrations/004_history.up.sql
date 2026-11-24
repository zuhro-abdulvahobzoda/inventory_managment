CREATE TABLE IF NOT EXISTS history (
    id         SERIAL PRIMARY KEY,
    item_id    INT NOT NULL REFERENCES items(id),
    user_id    INT NOT NULL REFERENCES users(id),
    action     VARCHAR(20) NOT NULL,
    changed_at TIMESTAMP NOT NULL DEFAULT now(),
    old_data   JSONB,
    new_data   JSONB
);

CREATE INDEX IF NOT EXISTS idx_history_item ON history(item_id);
