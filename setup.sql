-- the database will have 3 tables

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL UNIQUE 
);

CREATE TABLE tasks (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	category_id INT REFERENCES categories(id),
	title VARCHAR(255) NOT NULL,
	description TEXT,
	deadline TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users(id),
	title TEXT NOT NULL
);