INSERT INTO teams (team_name) VALUES 
('backend'),
('frontend'),
('payments'),
('mobile');

INSERT INTO users (user_id, username, team_name, is_active) VALUES 
('u1', 'Alice', 'backend', true),
('u2', 'Bob', 'backend', true),
('u3', 'Charlie', 'backend', true),
('u4', 'David', 'frontend', true),
('u5', 'Eva', 'frontend', true),
('u6', 'Frank', 'payments', true),
('u7', 'Grace', 'mobile', false);

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status) VALUES 
('pr-1005', 'Add search', 'u1', 'OPEN'),
('pr-1006', 'Fix authentication bug', 'u2', 'OPEN'),
('pr-1007', 'Update documentation', 'u4', 'MERGED'),
('pr-1008', 'Refactor API', 'u1', 'OPEN');

INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES 
('pr-1005', 'u2'),
('pr-1005', 'u3'),
('pr-1006', 'u1'),
('pr-1006', 'u3'),
('pr-1007', 'u5'),
('pr-1008', 'u2');