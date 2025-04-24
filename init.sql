-- Таблица пользователей (если нужно хранить UserData)
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    avatar TEXT,
    last_visit TIMESTAMP
);

-- Таблица постов
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    user_id UUID NOT NULL,
    user_name TEXT NOT NULL,
    user_avatar TEXT,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Таблица комментариев
CREATE TABLE comments (
    id UUID PRIMARY KEY,
    post_id UUID NOT NULL,
    user_id UUID NOT NULL,
    user_name TEXT NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
