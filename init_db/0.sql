-- Таблица для хранения информации о заказах
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    entry VARCHAR(50) NOT NULL,
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shard_key VARCHAR(10),
    sm_id INTEGER,
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    oof_shard VARCHAR(10)
);

-- Таблица для хранения информации о доставке
CREATE TABLE delivery (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    zip VARCHAR(20),
    city VARCHAR(100) NOT NULL,
    address VARCHAR(255) NOT NULL,
    region VARCHAR(100) NOT NULL,
    email VARCHAR(100)
);

-- Таблица для хранения информации об оплате
CREATE TABLE payment (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    transaction VARCHAR(255) NOT NULL,
    request_id VARCHAR(255),
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(50) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    payment_date_time BIGINT NOT NULL, -- Unix timestamp
    bank VARCHAR(100),
    delivery_cost DECIMAL(10, 2),
    goods_total DECIMAL(10, 2),
    custom_fee DECIMAL(10, 2)
);

-- Таблица для хранения информации о товарах
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    chrt_id INTEGER NOT NULL,
    track_number VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    rid UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(50),
    total_price DECIMAL(10, 2) NOT NULL,
    nm_id INTEGER NOT NULL,
    brand VARCHAR(100),
    status INTEGER NOT NULL
);