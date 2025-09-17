-- Mixed SQL Examples for Testing
-- Contains both good and bad practices to test comprehensive rule checking

-- GOOD: Well-designed table with proper naming and constraints
CREATE TABLE customer_accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_number VARCHAR(20) NOT NULL UNIQUE COMMENT 'Unique account identifier',
    customer_name VARCHAR(100) NOT NULL COMMENT 'Full customer name',
    email VARCHAR(255) NOT NULL COMMENT 'Customer email address',
    phone VARCHAR(20) COMMENT 'Contact phone number',
    account_type ENUM('personal', 'business') NOT NULL DEFAULT 'personal',
    status ENUM('active', 'inactive', 'suspended') NOT NULL DEFAULT 'active',
    balance DECIMAL(15,2) NOT NULL DEFAULT 0.00 COMMENT 'Account balance',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Customer account information';

-- GOOD: Proper index naming
CREATE INDEX idx_customer_accounts_email ON customer_accounts(email);
CREATE INDEX idx_customer_accounts_account_number ON customer_accounts(account_number);
CREATE INDEX idx_customer_accounts_status_type ON customer_accounts(status, account_type);

-- BAD: Table without primary key
CREATE TABLE audit_logs (
    log_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- BAD: Poor column naming
CREATE TABLE temp_data (
    id INT PRIMARY KEY,
    d1 VARCHAR(50),
    d2 INT,
    d3 DATETIME
);

-- GOOD: Insert with specified columns
INSERT INTO customer_accounts (account_number, customer_name, email, phone, account_type)
VALUES
    ('ACC001', 'John Smith', 'john.smith@email.com', '+1234567890', 'personal'),
    ('ACC002', 'Jane Doe', 'jane.doe@email.com', '+0987654321', 'business');

-- BAD: Insert without specifying columns
INSERT INTO temp_data VALUES (1, 'test', 123, NOW());

-- GOOD: Select with proper WHERE clause and specific columns
SELECT id, account_number, customer_name, email, balance
FROM customer_accounts
WHERE status = 'active'
    AND account_type = 'business'
    AND balance > 1000.00
ORDER BY balance DESC
LIMIT 50;

-- BAD: Select all columns without WHERE clause
SELECT * FROM customer_accounts;

-- BAD: Dangerous UPDATE without proper WHERE clause
UPDATE customer_accounts SET status = 'active';

-- GOOD: Safe UPDATE with specific WHERE clause
UPDATE customer_accounts
SET status = 'inactive', updated_at = NOW()
WHERE balance = 0.00
    AND status = 'active'
    AND updated_at < DATE_SUB(NOW(), INTERVAL 365 DAY);

-- BAD: Leading wildcard LIKE query (performance issue)
SELECT * FROM customer_accounts
WHERE customer_name LIKE '%Smith%';

-- GOOD: Proper LIKE query without leading wildcard
SELECT id, customer_name, email
FROM customer_accounts
WHERE customer_name LIKE 'John%'
    AND status = 'active';

-- MIXED: Table with some good and some bad practices
CREATE TABLE transaction_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- GOOD: Has primary key
    account_id BIGINT NOT NULL,  -- GOOD: NOT NULL constraint
    t_type VARCHAR(20),  -- BAD: Poor column naming, should be 'transaction_type'
    amount DECIMAL(15,2) NOT NULL,
    description TEXT,  -- POTENTIALLY BAD: TEXT type might be inefficient
    processed_by INT,  -- BAD: No foreign key constraint reference
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes LONGTEXT  -- BAD: LONGTEXT can cause performance issues
);

-- BAD: Index with poor naming
CREATE INDEX ix1 ON transaction_history(account_id);

-- GOOD: Index with proper naming
CREATE INDEX idx_transaction_history_created_date ON transaction_history(created_date);

-- BAD: Function in WHERE clause that prevents index usage
SELECT * FROM transaction_history
WHERE YEAR(created_date) = 2024;

-- GOOD: Proper date range query that can use indexes
SELECT id, account_id, amount, description
FROM transaction_history
WHERE created_date >= '2024-01-01'
    AND created_date < '2025-01-01'
    AND amount > 0;

-- BAD: DELETE without WHERE clause
DELETE FROM audit_logs;

-- GOOD: DELETE with specific conditions
DELETE FROM transaction_history
WHERE created_date < DATE_SUB(NOW(), INTERVAL 7 YEAR)
    AND amount = 0;