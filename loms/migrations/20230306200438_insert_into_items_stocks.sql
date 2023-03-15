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
DELETE FROM items_stocks WHERE sku = 1076963 AND warehouse_id = 1 AND count = 111;
DELETE FROM items_stocks WHERE sku = 1076963 AND warehouse_id = 2 AND count = 222;
DELETE FROM items_stocks WHERE sku = 1076963 AND warehouse_id = 3 AND count = 333;
DELETE FROM items_stocks WHERE sku = 1148162 AND warehouse_id = 1 AND count = 111;
DELETE FROM items_stocks WHERE sku = 1148162 AND warehouse_id = 2 AND count = 222;
DELETE FROM items_stocks WHERE sku = 1148162 AND warehouse_id = 3 AND count = 333;
-- +goose StatementEnd
