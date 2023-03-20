-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_items (
    id serial PRIMARY KEY, 
    orders_id int8,
    sku int8,
    count int4 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_items;
-- +goose StatementEnd
