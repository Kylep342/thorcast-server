CREATE OR REPLACE FUNCTION trigger_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS geocodex (
    city VARCHAR,
    state VARCHAR(2),
    lat NUMERIC(24, 8) CHECK (lat BETWEEN -90.0 AND 90.0),
    lng NUMERIC(24, 8) CHECK (lng BETWEEN -180.0 AND 180.0),
    requests INTEGER CHECK (requests > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (city, state)
)
;

ALTER TABLE geocodex OWNER TO thorcast;

DROP TRIGGER geocodex_update_timestamp;

CREATE TRIGGER geocodex_update_timestamp
BEFORE UPDATE ON geocodex
FOR EACH ROW
EXECUTE PROCEDURE trigger_update_timestamp();

