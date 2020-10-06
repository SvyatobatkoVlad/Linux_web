package main

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
	/*

			I am using windows so env varieble doesn't work!

			a.Initialize(
		        os.Getenv("APP_DB_USERNAME"),
		        os.Getenv("APP_DB_PASSWORD"),
				os.Getenv("APP_DB_NAME"))
	*/
	a := App{}
	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSLLMode)
	a.Run(":8080")
}
