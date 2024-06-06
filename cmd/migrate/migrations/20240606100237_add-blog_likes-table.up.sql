CREATE TABLE IF NOT EXISTS blog_likes (
    id SERIAL PRIMARY KEY,
    value INTEGER NOT NULL,
    blog_id INTEGER REFERENCES blogs(id),
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX blog_likes_blog_id_idx ON blog_likes(blog_id);
CREATE INDEX blog_likes_user_id_idx ON blog_likes(user_id);