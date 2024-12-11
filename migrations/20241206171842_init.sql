-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE documents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    mime VARCHAR(50),
    is_public BOOLEAN DEFAULT FALSE,
    is_file BOOLEAN DEFAULT TRUE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE grants (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) REFERENCES users(login) ON DELETE CASCADE,
    document_id INT REFERENCES documents(id) ON DELETE CASCADE,
    permission VARCHAR(20) DEFAULT 'read',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_documents_user_id ON documents(user_id);
CREATE INDEX idx_grants_login ON grants(login);
CREATE INDEX idx_grants_document_id ON grants(document_id);
CREATE UNIQUE INDEX idx_grants_unique ON grants(login, document_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS grants;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
