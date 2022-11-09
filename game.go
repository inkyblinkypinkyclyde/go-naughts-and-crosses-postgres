package main

import (
	"database/sql"
	"fmt"

	// "reflect"

	_ "github.com/lib/pq"
)

// Game vars
type game struct {
	game   bool
	player string
}

func getGrid(db *sql.DB) []string { // Get the grid from the database return grid
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

func taketurn(g game, db *sql.DB) {
	var input int
	printGrid(getGrid(db))
	fmt.Println("It's your turn!")
	fmt.Printf("Enter a number for player %v :", g.player)
	fmt.Scanln(&input)
	updateGrid(db, input, g.player)
	updateTurn(db, g)
	printGrid(getGrid(db))
}

func printGrid(grid []string) {
	fmt.Printf(" %v | %v | %v \n %v | %v | %v \n %v | %v | %v \n", grid[0], grid[1], grid[2], grid[3], grid[4], grid[5], grid[6], grid[7], grid[8])

}
func updateGrid(db *sql.DB, position int, ox string) {
	_, err := db.Exec("UPDATE grid SET value = $1 WHERE value = $2", ox, position)
	if err != nil {
		fmt.Println(err)
	}
}

func updateTurn(db *sql.DB, g game) {
	fmt.Println("Updating turn")
	if g.player == "x" {
		_, err := db.Exec("UPDATE turns SET value = 'o' WHERE id = 1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Turn updated to o")
	}
	if g.player == "o" {
		_, err := db.Exec("UPDATE turns SET value = 'x' WHERE id = 1")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Turn updated to x")
	}
	fmt.Println("Turn updated")
}

func turnCheck(db *sql.DB) string {
	var turn string
	rows, err := db.Query("SELECT * FROM turns")
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var id int
		var value string
		rows.Scan(&id, &value)
		turn = value
	}
	return turn
}

func selectPlayer() string {
	var player string
	fmt.Println("Select player x or o")
	fmt.Scanln(&player)
	return player
}

func setupdb(db *sql.DB) {
	r1, err := db.Query("SELECT COUNT(*) FROM grid")
	if err != nil {
		fmt.Println(err)
	}
	var count int
	for r1.Next() {
		r1.Scan(&count)
	}
	if count == 0 {
		for i := 1; i < 10; i++ {
			_, err := db.Exec("INSERT INTO grid (value) VALUES ($1)", i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	r2, err := db.Query("SELECT COUNT(*) FROM turns")
	if err != nil {
		fmt.Println(err)
	}
	var count2 int
	for r2.Next() {
		r2.Scan(&count2)
	}
	if count2 == 0 {
		_, err := db.Exec("INSERT INTO turns (value) VALUES ('x')")
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Printf("added %v rows to grid and %v rows to turns", count, count2)

}

func resetGrid(db *sql.DB) {
	fmt.Println("Resetting grid")
	for i := 1; i < 10; i++ {
		_, err := db.Exec("UPDATE grid SET value = $1 WHERE id = $2", i, i)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Grid reset")
}

func startGame(g game, db *sql.DB) bool {
	setupdb(db)
	resetGrid(db)
	g.game = true
	return g.game
}

func main() {
	connStr := "postgresql://localhost/naughtsAndCrosses?user=richardgannon&password=postgres&sslmode=disable" // Our connection string
	db, err := sql.Open("postgres", connStr)                                                                   // Open a database connection
	if err != nil {                                                                                            // If there is an error
		fmt.Println(err)
	}

	var g game
	g.game = startGame(g, db)
	g.player = selectPlayer()

	for g.game == true {
		turn := turnCheck(db) // Get the turn
		if turn == g.player { // If it is the players turn
			taketurn(g, db)                      // Take a turn
			if winCheck(getGrid(db), g.player) { // If the player has won
				g.game = false // End the game
			}
		}

	}

}
