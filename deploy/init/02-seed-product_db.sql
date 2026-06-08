\c product_db;

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(255),
    category VARCHAR(100),
    price NUMERIC(10, 2),
    image_url VARCHAR(255),
    quantity INT DEFAULT 0
);

INSERT INTO products (product_name, category, price, image_url, quantity) VALUES
    ('iPhone 16', 'Electronics', 79999.00, 'https://example.com/iphone-16.jpg', 10),
    ('Samsung Galaxy S25', 'Electronics', 74999.00, 'https://example.com/samsung-galaxy-s25.jpg', 5),
    ('MacBook Air M4', 'Laptops', 114999.00, 'https://example.com/macbook-air-m4.jpg', 20),
    ('Nike Running Shoes', 'Footwear', 4999.00, 'https://example.com/nike-running-shoes.jpg', 50),
    ('Sony WH-1000XM6', 'Headphones', 29999.00, 'https://example.com/sony-wh-1000xm6.jpg', 15)
;
