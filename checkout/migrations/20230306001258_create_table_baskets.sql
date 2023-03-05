-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS baskets (
    user_id int8 PRIMARY KEY,
    sku int8,
    count int4
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS baskets;
-- +goose StatementEnd
