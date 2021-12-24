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