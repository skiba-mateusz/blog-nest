CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,
    title varchar(255) NOT NULL,
    content TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id),
    category_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX blogs_user_id_idx ON blogs(user_id);
CREATE INDEX blogs_category_id_idx ON blogs(category_id);
CREATE INDEX blogs_title_idx ON blogs(title);