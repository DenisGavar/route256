-- +goose Up
-- +goose StatementBegin
INSERT INTO items_stocks (sku, warehouse_id, count)
VALUES
(1625903, 4, 100),
(2618151, 4, 100),
(2956315, 4, 100),
(2958025, 4, 100),
(3596599, 4, 100),
(3618852, 4, 100),
(4288068, 4, 100),
(4465995, 4, 100)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM items_stocks WHERE sku = 1625903 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 2618151 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 2956315 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 2958025 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 3596599 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 3618852 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4288068 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4465995 AND warehouse_id = 4 AND count = 100;
-- +goose StatementEnd
