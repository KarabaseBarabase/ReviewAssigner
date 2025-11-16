SELECT 'Initializing database' AS status;

SELECT 'CREATE DATABASE review_service'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'review_service')\gexec

\c review_service;

SELECT 'Database review_service is ready' AS status;