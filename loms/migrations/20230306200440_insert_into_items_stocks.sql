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
SELECT 'down SQL query';
-- +goose StatementEnd
