PRAGMA foreign_keys = ON;


CREATE TABLE IF NOT EXISTS users (
    id        TEXT PRIMARY KEY,           
    email     TEXT NOT NULL UNIQUE,
    username  TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL,              
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS categories (
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS posts (
    id         TEXT PRIMARY KEY,           
    user_id    TEXT NOT NULL,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post_categories (
    post_id     TEXT NOT NULL,
    category_id INTEGER NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS comments (
    id         TEXT PRIMARY KEY,           
    post_id    TEXT NOT NULL,
    user_id    TEXT NOT NULL,
    content    TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS post_likes (
    user_id TEXT NOT NULL,
    post_id TEXT NOT NULL,
    value   INTEGER NOT NULL CHECK (value IN (1, -1)),   
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS comment_likes (
    user_id    TEXT NOT NULL,
    comment_id TEXT NOT NULL,
    value      INTEGER NOT NULL CHECK (value IN (1, -1)),
    PRIMARY KEY (user_id, comment_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS sessions (
    id         TEXT PRIMARY KEY,           
    user_id    TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


INSERT OR IGNORE INTO categories (name, slug) VALUES
    ('Général', 'general'),
    ('Technologie', 'tech'),
    ('Humour', 'humour'),
    ('Questions', 'questions'),
    ('Annonces', 'annonces');
