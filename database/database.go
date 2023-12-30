package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	fmt.Println("Postgres has been connected ")
	//initialzing the sql db

	// username := "akash"
	// password := "your_password"
	// dbName := "your_database"
	// dbHost := "localhost"
	// dbPort := 3306 //

	db, err := sql.Open("postgres", "user:password@tcp(127.0.0.1:5432)/testdb")
	if err != nil {
		panic(err)
	}
	fmt.Println("Postgres has been connected ")

}

func main() {
	defer db.Close()
}
