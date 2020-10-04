package main

import (
	"fmt"
	"os"
)

var (
	dbUser     = os.Getenv("APP_DB_USERNAME")
	dbPassword = os.Getenv("APP_DB_PASSWORD")
	dbName     = os.Getenv("APP_DB_NAME")
	dbHost     = os.Getenv("APP_DB_HOST")
	dbPort     = os.Getenv("APP_DB_PORT")
	dbSLLMode  = os.Getenv("DB_SSLMODE")
)

func main() {
	dbUser = os.Getenv("APP_DB_USERNAME")
	fmt.Println("dbUser = ", dbUser)
	dbPassword = os.Getenv("APP_DB_PASSWORD")
	fmt.Println("APP_DB_PASSWORD = ", dbPassword)
	dbName = os.Getenv("APP_DB_NAME")
	fmt.Println("APP_DB_NAME = ", dbName)
	dbHost = os.Getenv("APP_DB_HOST")
	fmt.Println("APP_DB_HOST = ", dbHost)
	dbPort = os.Getenv("APP_DB_PORT")
	fmt.Println("APP_DB_PORT = ", dbPort)
	dbSLLMode = os.Getenv("DB_SSLMODE")
	fmt.Println("DB_SSLMODE = ", dbSLLMode)

	a := App{}

	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSLLMode)

	a.Run(":8080")
}
