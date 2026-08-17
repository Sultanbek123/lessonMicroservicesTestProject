CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    order_id INTEGER NOT NULL,
    message VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'created'
);
