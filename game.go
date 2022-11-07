package main

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"
)

func getGrid(db *sql.DB) []string {
	var grid []string
	rows, err := db.Query("SELECT * FROM grid ORDER BY id ASC")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var id int
		var value string
		rows.Scan(&id, &value)
		grid = append(grid, value)
	}
	fmt.Println(reflect.TypeOf(grid))
	return grid
}

func winCheck(grid []string, ox string) bool {
	if grid[0] == ox && grid[1] == ox && grid[2] == ox {
		return true
	}
	if grid[3] == ox && grid[4] == ox && grid[5] == ox {
		return true
	}
	if grid[6] == ox && grid[7] == ox && grid[8] == ox {
		return true
	}
	if grid[0] == ox && grid[3] == ox && grid[6] == ox {
		return true
	}
	if grid[1] == ox && grid[4] == ox && grid[7] == ox {
		return true
	}
	if grid[2] == ox && grid[5] == ox && grid[8] == ox {
		return true
	}
	if grid[0] == ox && grid[4] == ox && grid[8] == ox {
		return true
	}
	if grid[2] == ox && grid[4] == ox && grid[6] == ox {
		return true
	}
	return false
}

func printGrid(db *sql.DB) {
	grid := getGrid(db)
	if winCheck(grid, "O") {
		fmt.Println("O wins")
	}
	if winCheck(grid, "X") {
		fmt.Println("X wins")
	}
	fmt.Printf(" %v | %v | %v \n %v | %v | %v \n %v | %v | %v \n", grid[0], grid[1], grid[2], grid[3], grid[4], grid[5], grid[6], grid[7], grid[8])

}

func update(db *sql.DB, position int, ox string) {
	_, err := db.Exec("UPDATE grid SET value = $1 WHERE value = $2", ox, position)
	if err != nil {
		fmt.Println(err)
	}

}

func setup(db *sql.DB) { // Setup the database
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

func startGame(db *sql.DB) { // Start the game
	for i := 1; i < 10; i++ {
		_, err := db.Exec("UPDATE grid SET value = $1 WHERE id = $2", i, i)
		if err != nil {
			fmt.Println(err)
		}
	}
	_, err := db.Exec("UPDATE turns SET value = 1 WHERE id = 1")
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
	// setup(db)
	startGame(db)
	printGrid(db)

	turn := 0
	for true {
		input := 0
		var ox string
		if turn%2 == 0 {
			ox = "O"
		} else {
			ox = "X"
		}
		fmt.Printf("Enter a number for %v :", ox)
		fmt.Scanln(&input)
		update(db, input, ox)
		printGrid(db)
		turn += 1
	}

}
