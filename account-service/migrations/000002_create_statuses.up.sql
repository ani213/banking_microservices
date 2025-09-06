CREATE TABLE IF NOT EXISTS account_status (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,      -- e.g., "active", "inactive", "frozen"
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Insert default account_status
INSERT INTO account_status (name, description) VALUES
('active', 'Account is active and usable'),
('inactive', 'Account is inactive, not in use'),
('frozen', 'Account is frozen due to compliance or fraud'),
('closed', 'Account has been closed by the customer');
