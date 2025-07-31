CREATE TABLE datamodel (
        id BIGSERIAL PRIMARY KEY,
        name VARCHAR(255) NOT NULL
);

CREATE TABLE version (
        id BIGSERIAL PRIMARY KEY,
        object_type VARCHAR(255) NOT NULL,
        object_id BIGINT NOT NULL,
        json JSONB NOT NULL,
        version INT NOT NULL,
        action VARCHAR(255) NOT NULL,
        actor VARCHAR(255) NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users (
        id BIGSERIAL PRIMARY KEY,
        email VARCHAR(255) NOT NULL,
        status VARCHAR(255) NOT NULL,
        role VARCHAR(255),
        created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);