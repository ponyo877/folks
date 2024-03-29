CREATE DATABASE IF NOT EXISTS folks;

CREATE TABLE folks.rooms (
	id           VARCHAR(36)  NOT NULL,
	display_name VARCHAR(20)  NOT NULL,
    updated_at   TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    created_at   TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6) NOT NULL,
    PRIMARY KEY (id)
);