--orders schema
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


--warehouse schema
CREATE DATABASE warehousedb;

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

--payment
CREATE DATABASE paymentdb;

CREATE TABLE IF NOT EXISTS payments (
    id serial primary key,
    user_id int,
    order_id int,
    status text NOT NULL DEFAULT 'NEW' CHECK (status IN ('NEW', 'COMPLETED', 'FAILED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
