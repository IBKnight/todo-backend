CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY UNIQUE,
    name          VARCHAR(255) NOT NULL,
    username      VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_list
(
    id          SERIAL PRIMARY KEY UNIQUE,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS todo_item
(
    id          SERIAL PRIMARY KEY UNIQUE,
    title       VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    done        BOOLEAN      NOT NULL DEFAULT false
);


CREATE TABLE IF NOT EXISTS users_list
(
    id      SERIAL PRIMARY KEY UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    list_id INTEGER NOT NULL REFERENCES todo_list (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS list_items
(
    id      SERIAL PRIMARY KEY UNIQUE,
    list_id INTEGER NOT NULL REFERENCES todo_list (id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES todo_item (id) ON DELETE CASCADE
);

CREATE INDEX ON users_list (user_id);
CREATE INDEX ON users_list (list_id);
CREATE INDEX ON list_items (list_id);