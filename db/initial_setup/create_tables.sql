-- the database will have 2 tables

CREATE TABLE public.categories (
	id SERIAL PRIMARY KEY,
	title TEXT NOT NULL
);

CREATE TABLE public.tasks (
	id SERIAL PRIMARY KEY,
	category_id INT REFERENCES categories(id),
	title VARCHAR(255) NOT NULL,
	description TEXT,
	deadline TIMESTAMP,
	completed BOOLEAN,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE public.users (
	id SERIAL PRIMARY KEY,
	username TEXT UNIQUE NOT NULL,
	password TEXT
);

-- creating a new user
-- INSERT INTO users (email, password) VALUES (
--   'johndoe@mail.com',
--   crypt('johnspassword', gen_salt('bf'))
-- );

-- verifying a new user

-- SELECT id 
--   FROM users
--  WHERE email = 'johndoe@mail.com' 
--    AND password = crypt('johnspassword', password);

-- SELECT id 
--   FROM users
--  WHERE email = 'johndoe@mail.com' 
--    AND password = crypt('wrongpassword', password);