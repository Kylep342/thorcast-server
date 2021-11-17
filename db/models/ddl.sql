CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION trigger_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE IF NOT EXISTS geocodex (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    city VARCHAR,
    state VARCHAR(2),
    lat NUMERIC(24, 8) NOT NULL CHECK (lat BETWEEN -90.0 AND 90.0),
    lng NUMERIC(24, 8) NOT NULL CHECK (lng BETWEEN -180.0 AND 180.0),
    requests INTEGER CHECK (requests > 0),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (city, state)
)
;

ALTER TABLE geocodex OWNER TO thorcast;

BEGIN;
DROP TRIGGER IF EXISTS geocodex_update_timestamp ON geocodex;

CREATE TRIGGER geocodex_update_timestamp
BEFORE UPDATE ON geocodex
FOR EACH ROW
EXECUTE PROCEDURE trigger_update_timestamp();
COMMIT;

CREATE TABLE IF NOT EXISTS states (
    id uuid NOT NULL DEFAULT uuid_generate_v4(),
    name varchar(20),
    code varchar(2),
    PRIMARY KEY (code)
)
;

INSERT INTO states (name, code) VALUES
    ('Alabama', 'AL'),
    ('Alaska', 'AK'),
    ('Arizona', 'AZ'),
    ('Arkansas', 'AR'),
    ('California', 'CA'),
    ('Colorado', 'CO'),
    ('Connecticut', 'CT'),
    ('Delaware', 'DE'),
    ('Florida', 'FL'),
    ('Georgia', 'GA'),
    ('Hawaii', 'HI'),
    ('Idaho', 'ID'),
    ('Illinois', 'IL'),
    ('Indiana', 'IN'),
    ('Iowa', 'IA'),
    ('Kansas', 'KS'),
    ('Kentucky', 'KY'),
    ('Louisiana', 'LA'),
    ('Maine', 'ME'),
    ('Maryland', 'MD'),
    ('Massachusetts', 'MA'),
    ('Michigan', 'MI'),
    ('Minnesota', 'MN'),
    ('Mississippi', 'MS'),
    ('Missouri', 'MO'),
    ('Montana', 'MT'),
    ('Nebraska', 'NE'),
    ('Nevada', 'NV'),
    ('New Hampshire', 'NH'),
    ('New Jersey', 'NJ'),
    ('New Mexico', 'NM'),
    ('New York', 'NY'),
    ('North Carolina', 'NC'),
    ('North Dakota', 'ND'),
    ('Ohio', 'OH'),
    ('Oklahoma', 'OK'),
    ('Oregon', 'OR'),
    ('Pennsylvania', 'PA'),
    ('Rhode Island', 'RI'),
    ('South Carolina', 'SC'),
    ('South Sakota', 'SD'),
    ('Tennessee', 'TN'),
    ('Texas', 'TX'),
    ('Utah', 'UT'),
    ('Vermont', 'VT'),
    ('Virginia', 'VA'),
    ('Washington', 'WA'),
    ('West Virginia', 'WV'),
    ('Wisconsin', 'WI'),
    ('Wyoming', 'WY')
;