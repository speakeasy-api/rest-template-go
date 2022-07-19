/* allows the use of the uuid type */
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

/* allows regex checks on email domain */
CREATE EXTENSION IF NOT EXISTS citext;

/* helps with creating partial indexes for LIKE queries */
CREATE EXTENSION IF NOT EXISTS pg_trgm;

/* attempts to validate emails using regex before being inserted to the database  */
CREATE DOMAIN email AS citext CHECK (
    value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$'
);

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4 (),
    first_name VARCHAR (255) NOT NULL,
    last_name VARCHAR (255) NOT NULL,
    nickname VARCHAR (255) NOT NULL CHECK (nickname <> ''),
    password VARCHAR (255) NOT NULL CHECK (password <> ''),
    email email NOT NULL,
    country VARCHAR (255) NOT NULL CHECK (country <> ''),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT id_unique UNIQUE (id),
    CONSTRAINT nickname_unique UNIQUE (nickname),
    CONSTRAINT email_unique UNIQUE (email)
);

/* first name index */
CREATE INDEX idx_users_first_name ON users (first_name);

CREATE INDEX trgm_idx_users_first_name ON users USING gin (first_name gin_trgm_ops);

/* last name index */
CREATE INDEX idx_users_last_name ON users (last_name);

CREATE INDEX trgm_idx_users_last_name ON users USING gin (last_name gin_trgm_ops);

/* nickname index */
CREATE UNIQUE INDEX idx_users_nickname ON users (nickname);

CREATE INDEX trgm_idx_users_nickname ON users USING gin (nickname gin_trgm_ops);

/* email index */
CREATE UNIQUE INDEX email_idx ON users (email);

CREATE INDEX trgm_idx_users_email ON users USING gin (email gin_trgm_ops);

/* country index */
CREATE INDEX idx_users_country ON users (country);

CREATE INDEX trgm_idx_users_country ON users USING gin (country gin_trgm_ops);