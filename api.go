package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"os"
	"github.com/joho/godotenv"
	"github.com/emvi/null"
)

// structs
type Task struct {
	Id int // note: make sure the attributes are Capitalized -> if not they won't be exported -> json-encoder will not be able to access attributes, causing an empty object ("{}") to be returned.
	Title string
	Description string
	Category string
	Deadline null.Time
	Created_at null.Time
	Updated_at null.Time
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
	Category_Id string `json:"category_id"`
	Deadline null.Time `json:"deadline"`	
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware());

	/* --------------------------------------------------------------- URL ENDPOINTS -------------- */

	// ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello!")
	})

	// get all tasks
	r.GET("/alltasks", func(c *gin.Context) {
		var taskList []Task = getAllTasks()
		c.JSON(200, taskList)
	})

	// get a specific task by id
	r.POST("/gettask", func(c *gin.Context) {
		var id int;
		err := c.BindJSON(&id);
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", err)
			os.Exit(1)
		}

		var t Task = getTask(id);

		// fmt.Println(t);

		c.JSON(200, t)
	})

	// update a specific task by id
	r.POST("/updatetask", func(c *gin.Context) {
		var params UpdateTaskParams
		err := c.BindJSON(&params);
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", err)
			os.Exit(1)
		}

		updateTask(params)

		c.JSON(200, fmt.Sprintf("Successfully updated task with id: %v", params.Id))
	})

	// deletes a task
	r.POST("/deletetask", func(c *gin.Context) {
		var id int;
		err := c.BindJSON(&id);
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", err)
			os.Exit(1)
		}

		deleteTask(id);

		c.String(200, fmt.Sprintf("Successfully deleted task with id: %v", id))
	})

	// adds a task
	r.POST("/addtask", func(c *gin.Context) {
		var params CreateTaskParams
		err := c.BindJSON(&params)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to parse JSON body: %v\n", err)
			os.Exit(1)
		}

		// fmt.Println(params)
		addTask(params)		
	})

	// start the server
	r.Run()
}

/* ----------------------------------------------------------------- DATABASE FUNCTIONS --------- */
/* Initialises and returns a connection to the database */
func connectDB() (c *pgx.Conn) {
	// load the .env file that contains postgresql connection details
	godotenv.Load(".env")

	// open a connection to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// check that the connection is successfully established
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn;
}

/* Returns an array of Tasks stored in the database */
func getAllTasks() ([]Task) {
	c := connectDB()
	defer c.Close(context.Background())

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_all_tasks();")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch tasks from db: %v\n", err)
		os.Exit(1)
	}
	defer tasks.Close();

	var taskSlice []Task
	for tasks.Next() {
		var t Task
		err = tasks.Scan(
			&t.Id, 
			&t.Title,
			&t.Description,
			&t.Category,
			&t.Deadline,
			&t.Created_at,
			&t.Updated_at,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch tasks from db: %v\n", err)
			os.Exit(1)
		}
		taskSlice = append(taskSlice, t)
	}

	return taskSlice;
}

/* Return a Task by its id */
func getTask(id int) (Task) {
	c := connectDB()
	defer c.Close(context.Background())

	var t Task

	err := c.QueryRow(context.Background(), "SELECT * from public.get_all_tasks() WHERE id=$1;", id).Scan(
		&t.Id, 
		&t.Title,
		&t.Description,
		&t.Category,
		&t.Deadline,
		&t.Created_at,
		&t.Updated_at,		
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch task from db: %v\n", err)
		os.Exit(1)
	}

	return t;
}

/* Update a Task by its id */
func updateTask(t UpdateTaskParams) {
	c := connectDB()
	defer c.Close(context.Background())

	_, err := c.Exec(context.Background(), "UPDATE tasks SET category_id=$1, title=$2, description=$3, deadline=$4 WHERE id=$5", t.Category_Id, t.Title, t.Description, t.Deadline, t.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to update task: %v\n", err)
		os.Exit(1)
	}
}

/* Deletes a Task in the database with the corresponding id */
func deleteTask(id int) {
	c := connectDB()
	defer c.Close(context.Background())

	// use Exec to execute a query that does not return a result set
	commandTag, err := c.Exec(context.Background(), "DELETE FROM tasks where id=$1;", id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to delete task in db: %v\n", err)
		os.Exit(1)
	}
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "No row found to delete\n")
		os.Exit(1)
	}
}

/* Adds a Task to the database */
func addTask(params CreateTaskParams) {
	c := connectDB()
	defer c.Close(context.Background())

	commandTag, err := c.Exec(context.Background(), "INSERT INTO tasks (category_id, title, description, deadline) VALUES ($1, $2, $3, $4);", params.Category_Id, params.Title, params.Description, params.Deadline)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to add task to db: %v\n", err)
		os.Exit(1)
	}
	if commandTag.RowsAffected() != 1 {
		fmt.Fprintf(os.Stderr, "Task not added to db\n")
		os.Exit(1)
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

// deletes a task by its id (which is its primary-key in the db)
// 		curl -X POST 0.0.0.0:8080/deletetask -H "Content-Type: application/json" -d '2'
// 		curl -X POST https://tomato-backend-api.herokuapp.com/deletetask -H "Content-Type: application/json" -d '2'

// add a task
//		curl -X POST 0.0.0.0:8080/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": "2018-04-13T19:24:00+08:00"}'
//		curl -X POST https://tomato-backend-api.herokuapp.com/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": "2018-04-13T19:24:00+08:00"}'
//		curl -X POST 0.0.0.0:8080/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": null}'

// update a task
//		curl -X POST 0.0.0.0:8080/updatetask -H "Content-Type: application/json" -d '{"id":8, "category_id":"1", "title":"updated", "description":"this is an updated description", "deadline": "2018-04-13T19:24:00+08:00"}'














