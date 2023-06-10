-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sku_stocks (
    sku bigint,
    warehouse_id bigint,
    amount bigint,
    PRIMARY KEY (sku, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sku_stocks;
-- +goose StatementEnd