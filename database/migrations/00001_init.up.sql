CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT NOT NULL UNIQUE,
                       password TEXT NOT NULL,
                       role TEXT NOT NULL DEFAULT 'user'
);

CREATE TABLE manga (
                       id SERIAL PRIMARY KEY,
                       title TEXT NOT NULL,
                       description TEXT NOT NULL,
                       genre TEXT NOT NULL,
                       cover TEXT NOT NULL
);

CREATE TABLE comments (
                          id SERIAL PRIMARY KEY,
                          manga_id INT NOT NULL REFERENCES manga(id) ON DELETE CASCADE,
                          user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          text TEXT NOT NULL,
                          created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
