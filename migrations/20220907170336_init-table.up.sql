CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    host TEXT,
    owner TEXT,
    repository TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS scans (
    id SERIAL PRIMARY KEY,
    repository_id INT8 REFERENCES repositories(id),
    status TEXT,
    findings JSONB,
    queued_at TIMESTAMP DEFAULT NOW(),
    scanning_at TIMESTAMP,
    finished_at TIMESTAMP
);