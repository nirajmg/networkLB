package main
import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)
func main() {
	var (
		//password = os.Setenv("SQL_DB_PASSWORD")
		//user = os.Setenv("SQL_DB_USER")
		//port = os.Setenv("SQL_DB_PORT")
		//database = os.Setenv("MSSQL_DB_DATABASE")
		password = "root"
		user = "root"
		port = "5432:30788"
		database = "mydb"
	 )

	connectionString := fmt.Sprintf("user id=%s;password=%s;port=%s;database=%s", user, password, port, database)

	db, connectionError := sql.Open("postgres", connectionString); if connectionError != nil {
      fmt.Println(fmt.Errorf("error opening database: %v", connectionError))
   }

   err := db.Close()
   if err != nil {
	fmt.Printf("error: %v", err)
   }

}


