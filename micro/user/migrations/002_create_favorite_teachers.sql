-- Migration to create the favorite_teachers table
-- This table stores which teachers a user has marked as a favorite.

CREATE TABLE IF NOT EXISTS favorite_teachers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    teacher_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (user_id, teacher_id)
);