DROP TABLE IF EXISTS geocodex;

CREATE TABLE IF NOT EXISTS geocodex (
    city VARCHAR,
    state VARCHAR(2),
    lat NUMERIC(24, 8) CHECK (lat BETWEEN -90.0 AND 90.0),
    lng NUMERIC(24, 8) CHECK (lng BETWEEN -180.0 AND 180.0),
    requests INTEGER CHECK (requests > 0),
    created_at TIMESTAMP DEFAULT NOW()
)
;

ALTER TABLE geocodex ADD PRIMARY KEY (city, state);

ALTER TABLE geocodex OWNER TO thorcast;