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
SELECT 'down SQL query';
-- +goose StatementEnd
