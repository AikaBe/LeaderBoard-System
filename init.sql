-- Таблица пользователей (если нужно хранить UserData)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,  -- Используем SERIAL для автоинкремента
    name TEXT NOT NULL,
    avatar TEXT,
    last_visit TIMESTAMP
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    name TEXT,
    avatar TEXT,
    last_visit TIMESTAMP
);
-- Таблица постов
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,  -- Автоинкремент
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    user_name TEXT NOT NULL,  -- Имя пользователя
    user_avatar TEXT,         -- Аватар пользователя
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    archived_at TIMESTAMP,
    is_hidden BOOLEAN DEFAULT FALSE
);

-- Таблица комментариев
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,  -- Автоинкремент
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,  -- Внешний ключ на posts
    parent_comment_id INT REFERENCES comments(id) ON DELETE CASCADE,  -- Внешний ключ на comments
    user_name TEXT NOT NULL,  -- Имя пользователя
    user_avatar TEXT,         -- Аватар пользователя
    text TEXT NOT NULL,
    image_url TEXT,   
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
