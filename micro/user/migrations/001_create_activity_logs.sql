-- +migrate Up
-- This migration creates the activity_logs table to store
-- user activity entries. Each log is associated with a user and
-- captures the action performed along with an optional description
-- and the timestamp of when the activity occurred.

CREATE TABLE IF NOT EXISTS activity_logs (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS activity_logs;