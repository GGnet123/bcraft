CREATE TABLE users 
(
    id              SERIAL PRIMARY KEY NOT NULL,
    username        VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL
);

CREATE TABLE recipes (
    id          SERIAL PRIMARY KEY NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description VARCHAR(1000) NOT NULL,
    img_path        VARCHAR(500)
);

CREATE TABLE user_recipe_ratings (
    id          SERIAL PRIMARY KEY NOT NULL,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipe_id   INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    rating      INTEGER NOT NULL CHECK (rating between 1 and 5)
);

CREATE TABLE ingredients (
    id      SERIAL PRIMARY KEY NOT NULL,
    name    VARCHAR(255) NOT NULL
);

CREATE TABLE cooking_steps (
    id              SERIAL PRIMARY KEY NOT NULL,
    recipe_id       INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    step            INTEGER NOT NULL,
    description     VARCHAR(1000) NOT NULL,
    img_path        VARCHAR(500),
    cook_min_time   INTEGER NOT NULL DEFAULT 1,
    UNIQUE(recipe_id, step)
);

CREATE TABLE recipe_ingredients (
    id              SERIAL PRIMARY KEY NOT NULL,
    recipe_id       INTEGER NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id   INTEGER NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE
);

INSERT INTO ingredients
SELECT generate_series(1,100), md5(random()::text);

INSERT INTO recipes
SELECT generate_series(1,100), md5(random()::text), md5(random()::text);

INSERT INTO cooking_steps
SELECT generate_series(1,100), generate_series(1,100), 1, md5(random()::text), NULL, floor(random() * (1-11+1) + 11)::int;

INSERT INTO cooking_steps
SELECT generate_series(101,150), generate_series(1,50), 2, md5(random()::text), NULL, floor(random() * (1-11+1) + 11)::int;

INSERT INTO recipe_ingredients
SELECT generate_series(1,100), generate_series(1,100), floor(random() * (1-99+1) + 99)::int;