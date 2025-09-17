-- Good SQL Examples for Testing
-- These examples should pass all SQL review rules

-- Table with proper primary key and naming
CREATE TABLE user_profiles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_name VARCHAR(50) NOT NULL COMMENT 'User login name',
    email VARCHAR(100) NOT NULL COMMENT 'User email address',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation time',
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User profile information';

-- Index with proper naming convention
CREATE INDEX idx_user_profiles_email ON user_profiles(email);
CREATE INDEX idx_user_profiles_created_at ON user_profiles(created_at);

-- Insert with all columns specified
INSERT INTO user_profiles (user_name, email, created_at, updated_at)
VALUES ('john_doe', 'john@example.com', NOW(), NOW());

-- Select with proper WHERE clause
SELECT id, user_name, email, created_at
FROM user_profiles
WHERE email = 'john@example.com'
AND created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY);

-- Update with WHERE clause
UPDATE user_profiles
SET user_name = 'john_smith', updated_at = NOW()
WHERE id = 1;

-- Delete with WHERE clause
DELETE FROM user_profiles
WHERE created_at < DATE_SUB(NOW(), INTERVAL 365 DAY);

-- Another well-designed table
CREATE TABLE order_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL COMMENT 'Reference to orders table',
    product_id BIGINT NOT NULL COMMENT 'Reference to products table',
    quantity INT NOT NULL DEFAULT 1 COMMENT 'Item quantity',
    unit_price DECIMAL(10,2) NOT NULL COMMENT 'Price per unit',
    total_price DECIMAL(10,2) NOT NULL COMMENT 'Total price for this line item',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation time'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Order line items';

-- Composite index with proper naming
CREATE INDEX idx_order_items_order_product ON order_items(order_id, product_id);