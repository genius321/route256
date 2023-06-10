-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS carts (
    user_id bigint,
    sku bigint,
    amount bigint,
    PRIMARY KEY (user_id, sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carts;
-- +goose StatementEnd
