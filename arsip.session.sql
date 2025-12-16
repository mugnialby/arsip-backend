CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(128) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(128) NOT NULL,
    department_id INT,
    role_id INT,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    department_id INT,
    role_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    department_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE archive_characteristics (
    id SERIAL PRIMARY KEY,
    archive_characteristic_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE archive_types (
    id SERIAL PRIMARY KEY,
    archive_type_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE archive_hdr (
    id SERIAL PRIMARY KEY,
    archive_name VARCHAR(128) NOT NULL,
    archive_characteristic_id INT,
    archive_type_id INT,
    archive_date DATE,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE archive_attachments (
    id SERIAL PRIMARY KEY,
    archive_hdr_id INT NOT NULL,
    file_name VARCHAR(128),
    file_location VARCHAR(512),
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

drop table users;
drop table roles;
drop table archive_hdr;
drop table archive_attachments;
select * from users;
select * from roles;
select * from departments;

DELETE FROM USERS WHERE ID = 2;

truncate table roles;

INSERT INTO users(user_id, password_hash, full_name, department_id, role_id, STATUS, CREATED_BY, CREATED_AT)
VALUES('ALBYAM', 'ASD', 'ALBY AL MUGNI', 1, 1, 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

INSERT INTO roles(department_id, role_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES(1, 'SUPERADMIN','Y', 'ALBYAM', CURRENT_TIMESTAMP);

INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('PEMBINAAN', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('PIDANA UMUM', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('PIDANA KHUSUS', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('INTELIJEN', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('DATUN', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO departments(department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('PAPBB', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

select * from archive_hdr;
select * from archive_attachments;
select * from archive_types;
select * from archive_characteristics;
select * from file_types;

truncate table archive_hdr;
truncate table archive_attachments;

SET TIME ZONE 'GMT+7';



INSERT INTO archive_types(archive_type_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('Surat', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO archive_types(archive_type_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('Nota Dinas', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO archive_types(archive_type_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('Berkas Perkara', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

INSERT INTO archive_characteristics(archive_characteristic_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('Biasa', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
INSERT INTO archive_characteristics(archive_characteristic_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('Rahasia', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

INSERT INTO file_types(file_type_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES('', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

