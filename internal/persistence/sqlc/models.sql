CREATE TABLE users (
                       id character varying(64) PRIMARY KEY NOT NULL,
                       name character varying(64) NOT NULL,
                           email TEXT NOT NULL UNIQUE,
                       roles TEXT[] NOT NULL DEFAULT '{}',
                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);