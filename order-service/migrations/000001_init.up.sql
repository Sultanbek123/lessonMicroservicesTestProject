CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created'
);
