CREATE TABLE IF NOT EXISTS users(
    user_id         UUID  PRIMARY KEY,
    first_name      VARCHAR (50) NOT NULL,
    last_name       VARCHAR (50) NOT NULL,
    contact         VARCHAR (300) NOT NULL,
    phone           VARCHAR (50) UNIQUE NOT NULL,
    email           VARCHAR (300) UNIQUE NOT NULL,
    password        VARCHAR (50) NOT NULL,
    date_birthday   DATE NOT NULL,
    gender          smallint NOT NULL,
    role            VARCHAR (30) NOT NULL,
    activated       boolean DEFAULT false,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER trg_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();