\connect postgres;

-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

-- Insert default roles
INSERT INTO roles (name) VALUES ('user'), ('moderator'), ('admin') ON CONFLICT DO NOTHING;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role_id INT REFERENCES roles(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Set default role for existing users
UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user') WHERE role_id IS NULL;

-- Insert a user with the 'admin' role
INSERT INTO users (email, first_name, last_name, password, role_id)
VALUES (
    'admin@example.com',
    'Admin',
    'Super',
    '$2a$14$4Cxw5/NK2ARnNMcE8/jnSuo6vATld5cO1yxSuWXwniqgIJIa39I7a',  -- It's best to hash passwords before inserting them in a real application
    (SELECT id FROM roles WHERE name = 'admin')
);