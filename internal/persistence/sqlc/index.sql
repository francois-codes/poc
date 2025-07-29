CREATE INDEX idx_version_object_type_version_desc
    ON version (object_id, object_type, version DESC)
    INCLUDE (json);

-- careful with this, creates an index to verify unicity, so can slow everything down.
-- NOT EXECUTED FOR NOW
ALTER TABLE version
    ADD CONSTRAINT unique_object_version
        UNIQUE (object_type, object_id, version);