SELECT add_continuous_aggregate_policy(
    'price_changes_24h',
    start_offset => INTERVAL '3 days',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour'
);
