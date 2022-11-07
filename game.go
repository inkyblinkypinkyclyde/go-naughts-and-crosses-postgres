package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

func printGrid(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM grid")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var grid []string
	var res string
	for rows.Next() {
		rows.Scan(&res)
		grid = append(grid, res)
	}
	fmt.Println(reflect.TypeOf(grid))
	fmt.Println(reflect.TypeOf(grid[0]))
	fmt.Println(grid)
	fmt.Printf(" %v | %v | %v \n %v | %v | %v \n %v | %v | %v \n", grid[0], grid[1], grid[2], grid[3], grid[4], grid[5], grid[6], grid[7], grid[8])

}

func startGame(db *sql.DB) { // Start the game

	for i := 1; i < 10; i++ {
		_, err := db.Exec("INSERT INTO grid (value) VALUES ($1)", i)
		if err != nil {
			fmt.Println(err)
		}
	}
	_, err := db.Exec("INSERT INTO turns (value) VALUES ('1')")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	connStr := "postgresql://localhost/naughtsAndCrosses?user=richardgannon&password=postgres&sslmode=disable" // Our connection string
	db, err := sql.Open("postgres", connStr)                                                                   // Open a database connection
	if err != nil {                                                                                            // If there is an error
		fmt.Println(err)
	}

	startGame(db)
	printGrid(db)
}
