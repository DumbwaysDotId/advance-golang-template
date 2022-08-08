package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	// Call ProfileRoutes() and ProductRoutes() function here ...
}
