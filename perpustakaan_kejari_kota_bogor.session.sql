CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(128) NOT NULL,
    role_id INT,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    role_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    category_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE library_item_types (
    id SERIAL PRIMARY KEY,
    library_item_type_name VARCHAR(128) NOT NULL,
    library_item_type_code VARCHAR(16) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE publishers (
    id SERIAL PRIMARY KEY,
    publisher_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE racks (
    id SERIAL PRIMARY KEY,
    rack_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(256) NOT NULL,
    author VARCHAR(128),
    publisher_id INT,
    publish_location VARCHAR(128),
    published_year INT CHECK (published_year > 0),
    category_id INT,
    rack_id INT,
    rack_row INT,
    library_item_type_id INT,
    library_item_origin_id INT,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE book_attachments (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    file_name VARCHAR(256) NOT NULL,
    file_base64 TEXT NOT NULL,
    file_type_id INT NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE file_types (
    id SERIAL PRIMARY KEY,
    file_type_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE library_item_origins (
    id SERIAL PRIMARY KEY,
    library_item_origin_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

DROP TABLE books;
DROP TABLE categories;
DROP TABLE roles;
DROP TABLE users;
DROP TABLE racks;
DROP TABLE library_item_types;

SELECT * FROM roles;
SELECT * FROM users;
SELECT * FROM BOOKS;
SELECT * FROM RACKS;
SELECT * FROM CATEGORIES;
SELECT * FROM library_item_origins;
TRUNCATE TABLE BOOKS;
TRUNCATE TABLE book_attachments;
TRUNCATE TABLE racks;
TRUNCATE TABLE categories;
update RACKS SET ID = 3;

SELECT * FROM "books" WHERE status = 'Y' AND title LIKE '%HUKUM%' ORDER BY title asc

INSERT INTO USERS(USER_ID, PASSWORD_HASH, FULL_NAME, ROLE_ID, STATUS, CREATED_BY, CREATED_AT) 
VALUES('ALBYAM', 'ASD', 'ALBY AL MUGNI', 1, 'Y', 'ALBYAM', CURRENT_TIMESTAMP);