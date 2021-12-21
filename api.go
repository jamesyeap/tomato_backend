package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"context"
	"os"
	"github.com/lib/pq"
	// "encoding/json"
)

// structs
type Task struct {
	Id int
	Title string
	Description string
	Category string
	Deadline pq.NullTime
	Created_at pq.NullTime
	Updated_at pq.NullTime
}

func main() {
	r := gin.Default()

	/* ----------------------------------- URL ENDPOINTS ----------------------------------- */

	// ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello!")
	})

	// get all tasks
	r.GET("/alltasks", func(c *gin.Context) {
		var taskList []Task = getAllTasks(connectDB())
		fmt.Println(taskList)

		// var jsonData []byte
		// jsonData, _ = json.Marshal(taskList)

		// fmt.Println(string(jsonData))

		c.JSON(200, taskList)
	})

	// start the server
	r.Run(":8080")

}

/* ----------------------------------- DATABASE FUNCTIONS ----------------------------------- */
func connectDB() (c *pgx.Conn) {
	// postgresql connection details
	url := "postgres://msrwewroudbvot:f4e6c0a6f144fa28e13ef92503c9ac36f256ec8dce7ae4e0b56f4aa21b1e77a2@ec2-34-198-122-185.compute-1.amazonaws.com:5432/d2trus57r2q0ch"
	os.Setenv("DATABASE_URL", url);

	// open a connection to the database
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// check that the connection is successfully established
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return conn;
}

func getAllTasks(c *pgx.Conn) ([]Task) {
	// get all tasks

	tasks, err := c.Query(context.Background(), "SELECT * from public.get_all_tasks();")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to fetch tasks from db: %v\n", err)
		os.Exit(1)
	}

	defer tasks.Close();
	defer c.Close(context.Background());

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

		// fmt.Println(t.created_at.Time.String())
	}

	// fmt.Println(taskSlice);

	return taskSlice;
}















