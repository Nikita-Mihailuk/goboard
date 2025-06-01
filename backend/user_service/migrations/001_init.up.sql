CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    photo_url TEXT,
    role TEXT NOT NULL DEFAULT 'user'
);

INSERT INTO users (name, email, password_hash) VALUES
( 'Ильгиз', 'test@gmail.com', '$2a$10$8ErcFt7uXtYabCzxFNmn5.Wk1Xo.jD1.dymb6cZq3MqVQH9aJv.OG'),
( 'Ильнур', 'qwerty@gmail.com', '$2a$10$73Mw4PWu9/eSfW4GcRqIu.kD0kQOyTk1va0YRQsdbTwSdB3VQdTa.'),
( 'Никита', 'nikita@gmail.com', '$2a$10$vStk0YDIgWlUZFaN5FFfVuUIcKSTlMmR/R4tR6gqONJufNhpClXyS');
