-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_id bigserial PRIMARY KEY,
    status_name text default 'new',
    user_id bigint
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
