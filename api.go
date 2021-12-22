package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"os"
	"github.com/lib/pq"
	"time"
	"github.com/joho/godotenv"
)

// structs
type Task struct {
	Id int // note: make sure the attributes are Capitalized -> if not they won't be exported -> json-encoder will not be able to access attributes, causing an empty object ("{}") to be returned.
	Title string
	Description string
	Category string
	Deadline pq.NullTime
	Created_at pq.NullTime
	Updated_at pq.NullTime
}

type CreateTaskParams struct {
	Title string `json:"title"`
	Description string `json:"description"`
	Category_Id string `json:"category_id"`
	Deadline time.Time `json:"deadline"`
}

func main() {
	r := gin.Default()

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

	// start the server at 0.0.0.0:8080
	r.Run(":8080")

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

// get a task where id=1
//		curl -X POST 0.0.0.0:8080/gettask -H "Content-Type: application/json" -d '2'

// deletes a task by its id (which is its primary-key in the db)
// 		curl -X POST 0.0.0.0:8080/deletetask -H "Content-Type: application/json" -d '2'

// add a task
//		curl -X POST 0.0.0.0:8080/addtask -H "Content-Type: application/json" -d '{"category_id":"1", "title":"buy milk", "description":"muz be lactose-free lolz", "deadline": "2018-04-13T19:24:00+08:00"}'














