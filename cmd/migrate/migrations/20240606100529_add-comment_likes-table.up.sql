CREATE TABLE IF NOT EXISTS comment_likes (
    id SERIAL PRIMARY KEY,
    value INTEGER NOT NULL,
    comment_id INTEGER REFERENCES comments(id),
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX comment_likes_comment_id_idx ON comment_likes(comment_id);
CREATE INDEX comment_likes_user_id_idx ON comment_likes(user_id);