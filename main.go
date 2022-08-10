package main

import (
	"dumbmerch/database"
	"dumbmerch/pkg/mysql"
	"dumbmerch/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	// import "godotenv" here ...
)

func main() {

	// Init godotenv here ...

	// initial DB
	mysql.DatabaseInit()

	// run migration
	database.RunMigration()

	r := mux.NewRouter()

	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	// Initialization "uploads" folder to public here ...

	fmt.Println("server running localhost:5000")
	http.ListenAndServe("localhost:5000", r)
}
