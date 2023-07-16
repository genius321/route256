-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notifications (
    order_id bigserial,
    status_name text default 'new',
    user_id bigint,
    created_at timestamptz default now(),
    PRIMARY KEY (order_id, status_name)
);

CREATE INDEX idx_user_id ON notifications(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_user_id;
DROP TABLE IF EXISTS notifications;
-- +goose StatementEnd
