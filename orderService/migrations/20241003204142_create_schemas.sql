-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    goods jsonb,
    status text NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'COMPLETED', 'FAILED'))
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_type text NOT NULL,
    payload jsonb NOT NULL ,
    status text NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'COMPLETED', 'FAILED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
