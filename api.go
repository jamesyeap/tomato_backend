package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"os"
	"github.com/joho/godotenv"
	"github.com/emvi/null"
	"golang.org/x/crypto/bcrypt"
)

// structs

// note: make sure the attributes in the struct are Capitalized
//		if not they won't be exported and cannot be accessed (kinda like private in Java),
//		as such, .BindJson will not be able to access attributes, 
//		causing an empty object ("{}") to be returned

type Credentials struct {
	Username string `json:username`
	Password string `json:password`
}

type User struct {
	Id int `json:id`
	Username string `json:username`
}

type Task struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id int `json:"category_id"`
	Category string `json:"category"`
	Deadline null.Time `json:"deadline"`
	Completed bool `json:"completed"`
	Created_at null.Time `json:"created_at"`
	Updated_at null.Time `json:"updated_at"`
}

type Category struct {
	Id int `json:"category_id"`
	Title string `json:"category_title"`
}

type CreateTaskParams struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id string `json:"category_id"`
	Deadline null.Time `json:"deadline"`
}

type UpdateTaskParams struct {
	Id int `json:id`
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id int `json:"category_id"`
	Deadline null.Time `json:"deadline"`	
}

type GetTaskByIdParams struct {
	Id int `json:"id"`
}

type GetTaskByCategoryIdParams struct {
	Category_Id int `json:category_id`
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	r := gin.Default()

	// allow CORS
	r.Use(CORSMiddleware());

	/* --------------------------------------------------------------- URL ENDPOINTS -------------- */

	// ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello!")
	})

	// sign a new user up
	r.POST("/signup", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params Credentials;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, cancel, err);

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 8)

		details, err := signUp(params.Username, hashedPassword, c, cancel)

		if (err == nil) {
			c.JSON(200, details);
		}
	})

	// log in an existing user
	r.POST("/login", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params Credentials;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, cancel, err);

		fmt.Println(params);

		details, err := logIn(params.Username, params.Password, c, cancel)

		if (err == nil) {
			c.JSON(200, details);
		}
	})

	// get all tasks
	r.POST("/alltasks", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var taskList []Task = getAllTasks(c, cancel);
		c.JSON(200, taskList)
	})

	// get all completed tasks
	r.GET("/completedtasks", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var taskList []Task = getCompletedTasks(c, cancel);
		c.JSON(200, taskList)
	})

	// get all incomplete tasks
	r.GET("/incompletetasks", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var taskList []Task = getIncompleteTasks(c, cancel);
		c.JSON(200, taskList)
	})

	// get a specific task by id
	r.POST("/gettask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params GetTaskByIdParams;
		err := c.BindJSON(&params);
		assertJSONSuccess(c, cancel, err);

		var t Task = getTask(params.Id, c, cancel);

		c.JSON(200, t)
	})

	// get tasks by category id
	r.POST("/gettaskbycategoryid", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params GetTaskByCategoryIdParams
		err := c.BindJSON(&params);
		assertJSONSuccess(c, cancel, err);

		var taskList []Task = getTaskByCategoryId(params.Category_Id, c, cancel);

		c.JSON(200, taskList)
	})

	// update a specific task by id
	r.POST("/updatetask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params UpdateTaskParams
		err := c.BindJSON(&params);
		assertJSONSuccess(c, cancel, err);

		updateTask(params, c, cancel)

		c.JSON(200, fmt.Sprintf("Successfully updated task with id: %v", params.Id))
	})

	// mark a task as completed by id
	r.POST("/completetask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, cancel, err);

		completeTask(params.Id, c, cancel);

		c.JSON(200, fmt.Sprintf("Successfully completed task with id: %v", params.Id))
	})

	// mark a task as incomplete by id
	r.POST("/incompletetask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, cancel, err);

		incompleteTask(params.Id, c, cancel);

		c.JSON(200, fmt.Sprintf("Successfully marked task as incomplete with id: %v", params.Id))
	})

	// deletes a task by id
	r.POST("/deletetask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params GetTaskByIdParams;
		err := c.BindJSON(&params)
		assertJSONSuccess(c, cancel, err);

		deleteTask(params.Id, c, cancel);

		c.String(200, fmt.Sprintf("Successfully deleted task with id: %v", params.Id))
	})

	// adds a task
	r.POST("/addtask", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());

		var params CreateTaskParams
		err := c.BindJSON(&params)
		assertJSONSuccess(c, cancel, err);

		addTask(params, c, cancel)		
	})

	// gets a list of all categories
	r.GET("/allcategories", func(c *gin.Context) {
		_, cancel := context.WithCancel(context.Background());
		var categoryList []Category = getAllCategories(c, cancel);
		c.JSON(200, categoryList)
	})

	// start the server
	r.Run()
}

/* ----------------------------------------------------------------- DATABASE FUNCTIONS --------- */
/* Initialises and returns a connection to the database */
func connectDB(client *gin.Context, cancel context.CancelFunc) (c *pgx.Conn) {
	// load the .env file that contains postgresql connection details
	godotenv.Load(".env")

	// open a connection to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// check that the connection is successfully established
	assertDBSuccess(client, cancel, err);

	return conn;
}

/* Creates an account for a new user and returns the details of the user */
func signUp(username string, password []byte, client *gin.Context, cancel context.CancelFunc) (User, error) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2);", username, password);

	var user User

	// if a user with the same username already exists,
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to create a new user account: %v\n", err);

		// return HTTP Error 409: Conflict
		client.JSON(409, gin.H{"error": err.Error()});

		cancel();

		return user, err;
	}

	err = c.QueryRow(context.Background(), "SELECT id, username FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
	)
	assertDBOperationSuccess(client, cancel, err);

	return user, nil;
}

/* Logs in an existing-user */
func logIn(username string, password string, client *gin.Context, cancel context.CancelFunc) (User, error) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	var storedPassword string;
	var user User;

	err := c.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1;", username).Scan(&storedPassword);

	// check if the user exists
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to log user in: %v\n", err);

		// return HTTP Error 401: Unauthorised
		client.JSON(401, gin.H{"error": err.Error()});

		cancel();

		return user, err;
	}

	// check if the password is correct
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password));
	if (err != nil) {
		fmt.Fprintf(os.Stderr, "Unable to log user in: %v\n", err);

		client.JSON(401, gin.H{"error": err.Error()});

		cancel();

		return user, err;
	}

	// if credentials given are correct, return the user object
	err = c.QueryRow(context.Background(), "SELECT id, username FROM users WHERE username=$1", username).Scan(
		&user.Id,
		&user.Username,
	)
	assertDBOperationSuccess(client, cancel, err);

	return user, nil;
}

/* Returns an array of all Tasks stored in the database */
func getAllTasks(client *gin.Context, cancel context.CancelFunc) ([]Task) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_all_tasks();")
	assertDBOperationSuccess(client, cancel, err);
	defer tasks.Close();

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category_Id,
			&t.Category,
			&t.Deadline,
			&t.Completed,
			&t.Created_at,
			&t.Updated_at,	
		)
		assertDBOperationSuccess(client, cancel, err);
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

/* Returns an array of Tasks that belong to the specified ID stored in the database */
func getTaskByCategoryId(category_id int, client *gin.Context, cancel context.CancelFunc) ([]Task) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_tasks_in_category($1);", category_id)
	assertDBOperationSuccess(client, cancel, err);
	defer tasks.Close();

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category_Id,
			&t.Category,
			&t.Deadline,
			&t.Completed,
			&t.Created_at,
			&t.Updated_at,	
		)
		assertDBOperationSuccess(client, cancel, err);
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

func getCompletedTasks(client *gin.Context, cancel context.CancelFunc) ([]Task) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_completed_tasks();")
	assertDBOperationSuccess(client, cancel, err);
	defer tasks.Close();

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category_Id,
			&t.Category,
			&t.Deadline,
			&t.Completed,
			&t.Created_at,
			&t.Updated_at,	
		)
		assertDBOperationSuccess(client, cancel, err);
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

func getIncompleteTasks(client *gin.Context, cancel context.CancelFunc) ([]Task) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_incomplete_tasks();")
	assertDBOperationSuccess(client, cancel, err);
	defer tasks.Close();

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category_Id,
			&t.Category,
			&t.Deadline,
			&t.Completed,
			&t.Created_at,
			&t.Updated_at,	
		)
		assertDBOperationSuccess(client, cancel, err);
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

/* Return a Task by its id */
func getTask(id int, client *gin.Context, cancel context.CancelFunc) (Task) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	var t Task

	err := c.QueryRow(context.Background(), "SELECT * from public.get_all_tasks() WHERE id=$1;", id).Scan(
		&t.Id, 
		&t.Title,
		&t.Description,
		&t.Category,
		&t.Deadline,
		&t.Completed,
		&t.Created_at,
		&t.Updated_at,		
	)
	assertDBOperationSuccess(client, cancel, err);

	return t;
}

/* Update a Task by its id */
func updateTask(t UpdateTaskParams, client *gin.Context, cancel context.CancelFunc) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET category_id=$1, title=$2, description=$3, deadline=$4 WHERE id=$5;", t.Category_Id, t.Title, t.Description, t.Deadline, t.Id)
	assertDBOperationSuccess(client, cancel, err);
}

/* Mark a Task as completed by its id */
func completeTask(id int, client *gin.Context, cancel context.CancelFunc) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET completed='t' WHERE id=$1;", id);
	assertDBOperationSuccess(client, cancel, err);
}

/* Mark a previously completed task as incomplete by its id */
func incompleteTask(id int, client *gin.Context, cancel context.CancelFunc) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET completed='f' WHERE id=$1;", id);
	assertDBOperationSuccess(client, cancel, err);
}

/* Deletes a Task in the database with the corresponding id */
func deleteTask(id int, client *gin.Context, cancel context.CancelFunc) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	// use Exec to execute a query that does not return a result set
	commandTag, err := c.Exec(context.Background(), "DELETE FROM tasks where id=$1;", id)
	assertDBOperationSuccess(client, cancel, err);
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "No row found to delete\n")
		client.JSON(500, gin.H{"error": err.Error()})
		return;
	}
}

/* Adds a Task to the database */
func addTask(params CreateTaskParams, client *gin.Context, cancel context.CancelFunc) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	commandTag, err := c.Exec(context.Background(), "INSERT INTO tasks (category_id, title, description, deadline) VALUES ($1, $2, $3, $4);", params.Category_Id, params.Title, params.Description, params.Deadline)
	assertDBOperationSuccess(client, cancel, err);
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "Task not added to db\n")
		client.JSON(500, gin.H{"error": err.Error()})
		return;
	}
}

/* Returns a list of categories with their associated primary-keys */
func getAllCategories(client *gin.Context, cancel context.CancelFunc) ([]Category) {
	c := connectDB(client, cancel)
	defer c.Close(context.Background())

	categories, err := c.Query(context.Background(), "SELECT * from categories;")
	assertDBOperationSuccess(client, cancel, err);
	defer categories.Close();

	var categorySlice []Category
	for categories.Next() {
		var cat Category
		err = categories.Scan(
			&cat.Id,
			&cat.Title,
		)
		assertDBOperationSuccess(client, cancel, err);
		categorySlice = append(categorySlice, cat)
	}

	return categorySlice;
}

/* ------------------------------------------------------------ HELPER FUNCTIONS --------------------- */
// checks if there is an error connecting to the database,
//		if so, returns an error message to the client and cancels the context of the caller
func assertDBSuccess(client *gin.Context, cancel context.CancelFunc, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", e);

		// return http code of 500 to the client, which stands for "Internal Server Error"
		client.JSON(500, gin.H{"error": e.Error()});

		// halts execution of remaining functions to not do unnecessary work
		cancel();
	}
}

// checks if there is an error performing the specified request on the database,
//		if so, returns an error message to the client and cancels the context of the caller
func assertDBOperationSuccess(client *gin.Context, cancel context.CancelFunc, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to perform the requested action: %v\n", e);

		// return http code of 500 to the client, which stands for "Internal Server Error"
		client.JSON(500, gin.H{"error": e.Error()});

		// halts execution of remaining functions to not do unnecessary work
		cancel();
	}
}

// checks if there is an error connecting to the parsing JSON body,
//		if so, returns an error message to the client and stops execution of any remaining function-calls
func assertJSONSuccess(client *gin.Context, cancel context.CancelFunc, e error) {
	if (e != nil) {
		// print error message on server side so that its visible in the server logs
		fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", e);

		// return http code of 406 to the client, which stands for "Not Acceptable"
		client.JSON(406, gin.H{"error": e.Error()});

		// halts execution of remaining functions to not do unnecessary work
		cancel();
	}
}

/* ------ test-commands ------ */
// test if server is still up
// 		curl -X GET 0.0.0.0:8080/ping
//		curl -X GET https://tomato-backend-api.herokuapp.com/ping

// get all tasks
//		curl -X GET 0.0.0.0:8080/alltasks
//		curl -X GET https://tomato-backend-api.herokuapp.com/alltasks

// get a task where id=1
//		curl -X POST 0.0.0.0:8080/gettask -H "Content-Type: application/json" -d '2'
//		curl -X POST https://tomato-backend-api.herokuapp.com/gettask -H "Content-Type: application/json" -d '2'

// mark a task as complete with id
//		curl -X POST 0.0.0.0:8080/completetask -H "Content-Type: application/json" -d '2'

// mark a task as incomplete with id
//		curl -X POST 0.0.0.0:8080/incompletetask -H "Content-Type: application/json" -d '2'

// deletes a task by its id (which is its primary-key in the db)
// 		curl -X POST 0.0.0.0:8080/deletetask -H "Content-Type: application/json" -d '2'
// 		curl -X POST https://tomato-backend-api.herokuapp.com/deletetask -H "Content-Type: application/json" -d '2'

// add a task
//		curl -X POST 0.0.0.0:8080/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": "2018-04-13T19:24:00+08:00"}'
//		curl -X POST https://tomato-backend-api.herokuapp.com/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": "2018-04-13T19:24:00+08:00"}'
//		curl -X POST 0.0.0.0:8080/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": null}'

// update a task
//		curl -X POST 0.0.0.0:8080/updatetask -H "Content-Type: application/json" -d '{"id":8, "category_id":"1", "title":"updated", "description":"this is an updated description", "deadline": "2018-04-13T19:24:00+08:00"}'














