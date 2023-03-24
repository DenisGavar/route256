-- +goose Up
-- +goose StatementBegin
INSERT INTO items_stocks (sku, warehouse_id, count)
VALUES
(4487693, 4, 100),
(4669069, 4, 100),
(4678287, 4, 100),
(4678816, 4, 100),
(4679011, 4, 100),
(4687693, 4, 100),
(4996014, 4, 100),
(5097510, 4, 100)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM items_stocks WHERE sku = 4487693 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4669069 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4678287 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4678816 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4679011 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4687693 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 4996014 AND warehouse_id = 4 AND count = 100;
DELETE FROM items_stocks WHERE sku = 5097510 AND warehouse_id = 4 AND count = 100;
-- +goose StatementEnd
