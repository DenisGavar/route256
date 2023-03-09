-- +goose Up
-- +goose StatementBegin
INSERT INTO items_stocks (sku, warehouse_id, count)
VALUES
(1076963, 1, 111),
(1076963, 2, 222),
(1076963, 3, 333),
(1148162, 1, 111),
(1148162, 2, 222),
(1148162, 3, 333)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
