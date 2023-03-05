-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items_stocks_reservation (
    sku int8,
    warehouse_id int8,
    count int8,
    PRIMARY KEY (sku, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items_stocks_reservation;
-- +goose StatementEnd

