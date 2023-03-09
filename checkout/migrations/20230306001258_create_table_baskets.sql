-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS baskets (
    user_id int8,
    sku int8,
    count int4,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS baskets;
-- +goose StatementEnd
