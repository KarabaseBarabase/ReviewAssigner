-- users
CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    team_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- teams
CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(100) PRIMARY KEY,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- pull_requests
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    merged_at TIMESTAMPTZ NULL
);

-- pr_reviewers
CREATE TABLE IF NOT EXISTS pr_reviewers (
    id SERIAL PRIMARY KEY,
    pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    replaced_at TIMESTAMPTZ NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- Индексы
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_pr_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pr_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_pr_id ON pr_reviewers(pull_request_id);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
-- уникальный индекс для пар (pr, reviewer) когда запись активна
CREATE UNIQUE INDEX IF NOT EXISTS unique_active_reviewer ON pr_reviewers(pull_request_id, reviewer_id) WHERE is_active = true;


-- CREATE TABLE users (
--     user_id VARCHAR(50) PRIMARY KEY,
--     username VARCHAR(100) NOT NULL,
--     team_name VARCHAR(100) NOT NULL,
--     is_active BOOLEAN DEFAULT TRUE,
--     created_at TIMESTAMPTZ DEFAULT NOW(),
--     updated_at TIMESTAMPTZ DEFAULT NOW()
-- );

-- CREATE TABLE teams (
--     team_name VARCHAR(100) PRIMARY KEY,
--     created_at TIMESTAMPTZ DEFAULT NOW(),
--     updated_at TIMESTAMPTZ DEFAULT NOW()
-- );

-- CREATE TABLE pull_requests (
--     pull_request_id VARCHAR(50) PRIMARY KEY,
--     pull_request_name VARCHAR(255) NOT NULL,
--     author_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
--     status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
--     created_at TIMESTAMPTZ DEFAULT NOW(),
--     updated_at TIMESTAMPTZ DEFAULT NOW(),
--     merged_at TIMESTAMPTZ NULL
-- );

-- CREATE TABLE pr_reviewers (
--     id SERIAL PRIMARY KEY,
--     pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
--     reviewer_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
--     assigned_at TIMESTAMPTZ DEFAULT NOW(),
--     replaced_at TIMESTAMPTZ NULL,
--     is_active BOOLEAN DEFAULT TRUE
-- );

-- CREATE INDEX idx_users_team_active ON users(team_name, is_active);
-- CREATE INDEX idx_users_active ON users(is_active) WHERE is_active = true;
-- CREATE INDEX idx_pr_author_id ON pull_requests(author_id);
-- CREATE INDEX idx_pr_status ON pull_requests(status);
-- CREATE INDEX idx_pr_reviewers_pr_id ON pr_reviewers(pull_request_id);
-- CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers(reviewer_id);
-- CREATE UNIQUE INDEX unique_active_reviewer ON pr_reviewers(pull_request_id, reviewer_id) WHERE is_active = true;