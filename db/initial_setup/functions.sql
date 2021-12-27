-- get all tasks
CREATE OR REPLACE FUNCTION public.get_all_tasks()
	RETURNS TABLE 
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id;
END
$$;

-- get all completed tasks
CREATE OR REPLACE FUNCTION public.get_completed_tasks()
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			tasks.completed = 't';
END
$$;

-- get all outstanding tasks
CREATE OR REPLACE FUNCTION public.get_incomplete_tasks()
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			tasks.completed = 'f';
END
$$;

-- get tasks by category id
CREATE OR REPLACE FUNCTION public.get_tasks_in_category(Specified_Category_Id INT)
	RETURNS TABLE
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category_id INT,
			category TEXT,
			deadline TIMESTAMP,
			completed BOOLEAN,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)
	language plpgsql
AS
$$
BEGIN
	RETURN QUERY
		SELECT 
			tasks.id,
			tasks.title,
			tasks.description,
			categories.id,
			categories.title,
			tasks.deadline,
			tasks.completed,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id

		WHERE
			categories.id = Specified_Category_Id;
END
$$;
