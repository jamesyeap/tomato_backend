-- get all tasks
CREATE OR REPLACE FUNCTION public.get_all_tasks()
	RETURNS TABLE 
		(
			id INT,
			title VARCHAR(255),
			description TEXT,
			category TEXT,
			deadline TIMESTAMP,
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
			categories.title,
			tasks.deadline,
			tasks.created_at,
			tasks.updated_at
		FROM
			public.tasks
				INNER JOIN public.categories ON public.tasks.category_id=public.categories.id;
END
$$;