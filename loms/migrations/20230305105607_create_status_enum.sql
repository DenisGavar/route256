-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status as ENUM ('new', 'awaiting payment', 'failed', 'payed', 'cancelled');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE order_status;
-- +goose StatementEnd
