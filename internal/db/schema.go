package db

const schema = `
CREATE TABLE IF NOT EXISTS packages (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    summary    TEXT DEFAULT '',
    content    TEXT NOT NULL,
    tags       TEXT DEFAULT '[]',
    project    TEXT DEFAULT '',
    created_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_packages_name ON packages(name);
CREATE INDEX IF NOT EXISTS idx_packages_project ON packages(project);
CREATE INDEX IF NOT EXISTS idx_packages_expires_at ON packages(expires_at);
`
