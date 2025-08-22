-- Membuat tabel orders
CREATE TABLE IF NOT EXISTS orders (
    id CHAR(36) PRIMARY KEY,
    buyer_id CHAR(36) NOT NULL,
    order_date DATETIME NOT NULL,
    total DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    -- Relasi ke tabel users (buyer/pembeli)
    CONSTRAINT fk_orders_buyers FOREIGN KEY (buyer_id) REFERENCES users(id)
);
