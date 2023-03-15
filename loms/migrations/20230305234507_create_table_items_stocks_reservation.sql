-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items_stocks_reservation (
    id serial PRIMARY KEY, 
    sku int8,
    warehouse_id int8,
    orders_id int8,
    count int8 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items_stocks_reservation;
-- +goose StatementEnd

