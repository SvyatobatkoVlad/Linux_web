package main

import "fmt"

var (
	API_SECRET = "123532"
	dbUser     = "postgres"
	dbPassword = "root"
	dbName     = "web"
	dbHost     = "127.0.0.1"
	dbPort     = "5432"
	dbSLLMode  = "disable"
)

func main() {
	fmt.Println(dbUser, dbPassword, dbName, dbHost, dbPort, dbSLLMode)
	a := App{}

	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSLLMode)

	a.Run(":8080")
}
