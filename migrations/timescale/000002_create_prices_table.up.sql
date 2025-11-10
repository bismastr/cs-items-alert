CREATE TABLE prices (
    id BIGSERIAL,
    item_id INTEGER NOT NULL,
    sell_price INTEGER NOT NULL,
    sell_listings INTEGER NOT NULL,
    time TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

SELECT create_hypertable('prices', 'time');

CREATE INDEX idx_prices_item_id_time ON prices (item_id, time DESC);