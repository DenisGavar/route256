-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS outbox_orders (
    id serial PRIMARY KEY,
    orders_id int8,
    payload json,
    created_at timestamp,
    sent boolean
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox_orders;
-- +goose StatementEnd
