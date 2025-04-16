CREATE TYPE order_status AS ENUM ('PENDING', 'PROCESSING','COMPLETED', 'FAILED');

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    payment_address VARCHAR(255) NOT NULL,
    status order_status NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    order_expiration TIMESTAMP WITH TIME ZONE NOT NULL 
);

CREATE TABLE payment_addresses (
    address VARCHAR(255) PRIMARY KEY,
    private_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP 
);
