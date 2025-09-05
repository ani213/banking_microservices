CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,   -- unique user ID
    email TEXT UNIQUE NOT NULL,                      -- login email
    password_hash TEXT NOT NULL,                     -- hashed password
    full_name TEXT,                                  -- optional full name
    created_at TIMESTAMPTZ DEFAULT now(),            -- when user created
    updated_at TIMESTAMPTZ DEFAULT now()             -- last updated
);

-- Trigger to auto-update updated_at on row changes
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = now();
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();
