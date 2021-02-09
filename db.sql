CREATE DATABASE db
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'Russian_Russia.1251'
    LC_CTYPE = 'Russian_Russia.1251'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;
    
CREATE TABLE static
(
    date date NOT NULL,
    views integer DEFAULT 0,
    clicks integer DEFAULT 0,
    cost double precision DEFAULT 0.0
)

TABLESPACE pg_default;

ALTER TABLE static
    OWNER to postgres;