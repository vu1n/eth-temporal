package main

import (
	"database/sql"
	"eth-temporal/app"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", app.DbHost, app.DbPort, app.DbUser, app.DbPassword, app.DbName)

	fmt.Println(psqlconn)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to %s\n", app.DbHost)
	// clean up db connection
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS ethereum;")
	if err != nil {
		panic(err)
	}

	fmt.Println("Schema created")
}
