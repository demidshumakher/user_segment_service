-- Инициализация схемы
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS segments (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS user_segments (
    user_id INTEGER NOT NULL,
    segment TEXT NOT NULL,
    PRIMARY KEY (user_id, segment),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (segment) REFERENCES segments (name) ON DELETE CASCADE
);

-- Тестовые данные (100 пользователей)
INSERT INTO users (id)
SELECT gs
FROM generate_series(1, 100) AS gs
ON CONFLICT DO NOTHING;
