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
    archive_date DATE,
    archive_number VARCHAR(128),
    archive_name VARCHAR(128) NOT NULL,
    archive_characteristic_id INT,
    archive_type_id INT,
    department_id INT,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'archive_hdr';

CREATE INDEX ON archive_hdr(archive_date);
CREATE INDEX ON archive_hdr(archive_number);
CREATE INDEX ON archive_hdr(department_id);

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

CREATE INDEX ON archive_attachments(archive_hdr_id);

CREATE TABLE file_types (
    id SERIAL PRIMARY KEY,
    file_type_name VARCHAR(128) NOT NULL,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE TABLE archive_role_access (
    id SERIAL PRIMARY KEY,
    archive_hdr_id INT NOT NULL,
    role_id INT NOT NULL,
    department_id INT NOT NULL,
    -- can_view BOOLEAN DEFAULT TRUE,
    -- can_download BOOLEAN DEFAULT FALSE,
    -- can_edit BOOLEAN DEFAULT FALSE,
    status VARCHAR(1) DEFAULT 'Y' NOT NULL,
    created_by VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(128),
    modified_at TIMESTAMP
);

CREATE INDEX ON archive_role_access(archive_hdr_id);
CREATE INDEX ON archive_role_access(role_id);
CREATE INDEX ON archive_role_access(department_id);
CREATE INDEX ON archive_role_access(archive_hdr_id, role_id, department_id);

drop table users;
drop table roles;
drop table archive_hdr;
drop table archive_attachments;
drop table archive_role_access;

SELECT
    USR.USER_ID,
    USR.ROLE_ID,
    DEP.DEPARTMENT_NAME,
    ROL.ID ROL_ID,
    ROL.ROLE_NAME
FROM USERS USR
LEFT JOIN DEPARTMENTS DEP ON USR.department_id = DEP.ID
LEFT JOIN ROLES ROL ON USR.ROLE_id = ROL.ID
;

select * from users ORDER BY ID;
UPDATE USERS SET ID = 2 WHERE USER_ID = 'KAJARI'
DELETE FROM USERS WHERE ID > 1 AND ID NOT IN (14);
select * from roles ORDER BY ID;
select * from departments;

SELECT
    ROL.DEPARTMENT_ID,
    DEP.DEPARTMENT_NAME,
    ROL.ID ROLE_ID,
    ROL.ROLE_NAME
FROM ROLES ROL
LEFT JOIN departments DEP ON DEP.ID = ROL.department_id
;

UPDATE DEPARTMENTS SET ID = 2 WHERE ID = 1;
UPDATE roles SET DEPARTMENT_ID = 5 WHERE ID IN (18, 19);
select * from archive_hdr;
select * from archive_role_access;
select * from archive_attachments;

delete from users where id in (3, 4)

update USERS set PASSWORD_HASH = 'febyg' WHERE ID = 6;

DELETE FROM USERS WHERE ID = 2;
SELECT * FROM "roles" WHERE status = 'Y' AND roles.id NOT IN (1) ORDER BY role_name asc

truncate table roles;

INSERT INTO users(user_id, password_hash, full_name, department_id, role_id, STATUS, CREATED_BY, CREATED_AT)
VALUES('KAJARI', 'kajari', 'KAJARI', 1, 1, 'Y', 'ALBYAM', CURRENT_TIMESTAMP);

INSERT INTO roles(id, department_id, role_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES(2, 1, 'KAJARI','Y', 'ALBYAM', CURRENT_TIMESTAMP);

DELETE FROM ROLES WHERE ID = 27;
UPDATE ROLES SET MODIFIED_AT = NULL;

INSERT INTO departments(id, department_name, STATUS, CREATED_BY, CREATED_AT) 
VALUES(1, 'SUPERADMIN', 'Y', 'ALBYAM', CURRENT_TIMESTAMP);
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
select * from archive_role_access;
select * from archive_attachments;
select * from archive_types;
select * from archive_characteristics;
select * from file_types;

truncate table archive_hdr;
truncate table archive_attachments;
truncate table archive_role_access;

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
SELECT "archive_hdr"."id",
    "archive_hdr"."archive_name",
    "archive_hdr"."archive_number",
    "archive_hdr"."archive_characteristic_id",
    "archive_hdr"."archive_type_id",
    "archive_hdr"."archive_date",
    "archive_hdr"."department_id",
    "archive_hdr"."status",
    "archive_hdr"."created_by",
    "archive_hdr"."created_at",
    "archive_hdr"."modified_by",
    "archive_hdr"."modified_at"
FROM "archive_hdr"
    LEFT JOIN archive_role_access ON archive_role_access.archive_hdr_id = archive_hdr.id
    AND archive_role_access.role_id = 39
    AND archive_role_access.department_id = 6
WHERE archive_role_access.status = 'Y'
    AND archive_hdr.status = 'Y'
ORDER BY archive_date DESC,
    archive_name ASC