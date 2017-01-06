package app

import (
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-app/src/users/actions"
)

// SetupRoutes adds routes for this app to this router.
func SetupRoutes(r *router.Router) {

	// Resource Routes
	r.Add("/users", useractions.HandleIndex)
	r.Add("/users/create", useractions.HandleCreateShow)
	r.Add("/users/create", useractions.HandleCreate).Post()
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdateShow)
	r.Add("/users/{id:[0-9]+}/update", useractions.HandleUpdate).Post()
	r.Add("/users/{id:[0-9]+}/destroy", useractions.HandleDestroy).Post()
	r.Add("/users/{id:[0-9]+}", useractions.HandleShow)
	r.Add("/users/login", useractions.HandleLoginShow)
	r.Add("/users/login", useractions.HandleLogin).Post()
	r.Add("/users/logout", useractions.HandleLogout).Post()
	r.Add("/users/password", useractions.HandlePasswordReset)
	r.Add("/users/password/reset", useractions.HandlePasswordResetShow)
	r.Add("/users/password/reset", useractions.HandlePasswordResetSend).Post()
	r.Add("/users/password/sent", useractions.HandlePasswordResetSentShow)

	// Set the default file handler
	r.FileHandler = fileHandler
	r.ErrorHandler = errHandler

	// Add a files route to handle static images under files
	// - nginx deals with this in production - perhaps only do this in dev?
	r.Add("/files/{path:.*}", fileHandler)
	r.Add("/favicon.ico", fileHandler)

	// Add the home page route
	r.Add("/", homeHandler)

}
