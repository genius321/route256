-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sku_stocks_reservation (
    sku bigint,
    warehouse_id bigint,
    order_id bigint,
    amount bigint,
    PRIMARY KEY (sku, warehouse_id, order_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sku_stocks_reservation;
-- +goose StatementEnd