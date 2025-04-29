CREATE TABLE favorites (
                           id SERIAL PRIMARY KEY,
                           user_id INTEGER NOT NULL,
                           manga_id INTEGER NOT NULL,
                           CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                           CONSTRAINT fk_manga FOREIGN KEY (manga_id) REFERENCES mangas(id) ON DELETE CASCADE,
                           CONSTRAINT unique_favorite UNIQUE (user_id, manga_id)
);
