package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	rpcHost  = "https://eth-rpc.gateway.pokt.network"
	host     = "localhost"
	port     = 5433
	user     = "temporal"
	password = "temporal"
	dbname   = "postgres"
)

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Println(psqlconn)

	// Connect to pg db
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected to %s\n", host)
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
