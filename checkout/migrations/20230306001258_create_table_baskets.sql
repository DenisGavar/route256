-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS baskets (
    id serial PRIMARY KEY,
    user_id int8,
    sku int8,
    count int4 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS baskets;
-- +goose StatementEnd
