-- Bad SQL Examples for Testing
-- These examples should trigger various SQL review rules

-- Table without primary key (BAD)
CREATE TABLE bad_users (
    name VARCHAR(50),
    email VARCHAR(100)
);

-- Table with poor naming (BAD)
CREATE TABLE t1 (
    id INT AUTO_INCREMENT PRIMARY KEY,
    n VARCHAR(50),
    e VARCHAR(100)
);

-- Table with disallowed data types (BAD)
CREATE TABLE risky_table (
    id INT PRIMARY KEY,
    data TEXT,  -- Potentially problematic for large data
    blob_data LONGBLOB  -- Can cause performance issues
);

-- Insert without specifying columns (BAD)
INSERT INTO bad_users VALUES ('John', 'john@example.com');

-- Select without WHERE clause on large table (BAD)
SELECT * FROM bad_users;

-- Update without WHERE clause (DANGEROUS)
UPDATE bad_users SET name = 'Updated';

-- Delete without WHERE clause (DANGEROUS)
DELETE FROM bad_users;

-- Using SELECT * (BAD PRACTICE)
SELECT * FROM bad_users WHERE id = 1;

-- Table without charset specification (BAD)
CREATE TABLE no_charset (
    id INT PRIMARY KEY,
    name VARCHAR(50)
);

-- Index with poor naming (BAD)
CREATE INDEX ix1 ON bad_users(email);

-- Column without NOT NULL where it should be (BAD)
CREATE TABLE missing_constraints (
    id INT PRIMARY KEY,
    email VARCHAR(100),  -- Should be NOT NULL
    created_at TIMESTAMP  -- Should have default
);

-- Using reserved keywords as identifiers (BAD)
CREATE TABLE `order` (
    `select` INT PRIMARY KEY,
    `from` VARCHAR(50),
    `where` DATETIME
);

-- Dangerous functions in WHERE clause (BAD)
SELECT * FROM bad_users WHERE RAND() > 0.5;

-- Leading wildcard in LIKE (BAD for performance)
SELECT * FROM bad_users WHERE name LIKE '%john%';

-- Missing foreign key constraints (if policy requires them)
CREATE TABLE orphan_records (
    id INT PRIMARY KEY,
    user_id INT,  -- Should reference bad_users(id)
    content TEXT
);

-- Mixing DDL and DML in same statement block (BAD)
CREATE TABLE mixed_ops (id INT PRIMARY KEY);
INSERT INTO mixed_ops VALUES (1);

-- Column with excessive length (BAD)
CREATE TABLE excessive_lengths (
    id INT PRIMARY KEY,
    description VARCHAR(65535)  -- Excessive length
);