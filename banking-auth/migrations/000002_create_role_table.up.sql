CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,      
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Insert default account_status
INSERT INTO roles (name) VALUES
('customer'),
('teller'),
('manager'),
('loan_officer');
