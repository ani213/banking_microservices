CREATE TABLE IF NOT EXISTS account_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,         -- e.g., "savings", "current", "loan"
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Insert some default account types
INSERT INTO account_types (name, description) VALUES
('savings', 'Savings account for individuals'),
('current', 'Current account for businesses'),
('loan', 'Loan account');