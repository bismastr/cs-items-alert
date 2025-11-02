CREATE TABLE prices (
    id BIGSERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    sell_price INTEGER NOT NULL CHECK (sell_price >= 0),
    sell_listings INTEGER NOT NULL CHECK (sell_listings >= 0),
    created_at TIMESTAMPTZ DEFAULT NOW()
);