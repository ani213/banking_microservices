CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL, 
    display_name TEXT UNIQUE NOT NULL,     
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Insert default account_status
INSERT INTO roles (name,display_name) VALUES
('customer', 'Customer'),
('teller', 'Bank Teller'),
('manager', 'Branch Manager'),
('loan_officer', 'Loan Officer'),
('auditor', 'Auditor'),
('compliance', 'Compliance Officer'),
('admin', 'System Admin'),
('super_admin', 'Super Administrator');
