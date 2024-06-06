CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    blog_id INTEGER REFERENCES blogs(id),
    parent_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX comments_user_id_idx ON comments(user_id);
CREATE INDEX comments_blog_id_idx on comments(blog_id);
CREATE INDEX comments_parent_id_idx on comments(parent_id);