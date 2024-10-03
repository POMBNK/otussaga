-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    id serial primary key,
    name text NOT NULL,
    price int NOT NULL
);

CREATE TABLE IF NOT EXISTS order_reservations (
    id serial primary key,
    order_id int,
    goods jsonb,
    status text NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'COMPLETED', 'FAILED'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
