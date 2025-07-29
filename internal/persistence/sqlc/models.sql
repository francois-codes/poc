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
        created_at TIMESTAMPTZ NOT NULL,
        FOREIGN KEY (object_id) REFERENCES datamodel(id)
);
