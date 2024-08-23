package main

import "os"

// This assumes that you use environment variables APP_DB_USERNAME, APP_DB_PASSWORD, and APP_DB_NAME to store your database’s username, password, and name respectively.
func main() {
	a := App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	a.Run(":5432")
}

// We’re going to use PostgreSQL default parameters for the purposes of testing:
// export APP_DB_USERNAME=postgres
// export APP_DB_PASSWORD=root
// export APP_DB_NAME=postgres
// postgres://postgres:root@localhost:5432/export APP_DB_NAME=go_rest_api_db?sslmode=disable

// psql -h localhost -p 5432 -U postgres "dbname=go_rest_api_db sslmode=require"