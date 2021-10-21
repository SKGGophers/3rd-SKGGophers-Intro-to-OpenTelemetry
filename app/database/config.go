package database

import "fmt"

var (
	dbUsername = "postgres"
	dbPassword = "yolo123"
	dbHost     = "localhost"
	dbTable    = "postgres"
	dbPort     = "5432"
	pgConnStr  = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbHost, dbPort, dbUsername, dbTable, dbPassword)
)
