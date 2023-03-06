-- +goose Up
-- +goose StatementBegin
INSERT INTO items_stocks (sku, warehouse_id, count)
VALUES
(1, 11, 111),
(2, 11, 222),
(3, 11, 333),
(1, 22, 111),
(2, 22, 222),
(3, 22, 333),
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
