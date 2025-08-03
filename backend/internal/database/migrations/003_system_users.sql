-- +goose Up
CREATE TABLE system_users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    form_data JSONB,
    needed_min_btu DECIMAL,
    needed_max_btu DECIMAL
);

-- +goose Down
DROP TABLE system_users