package main

import (
	"fmt"

	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"
)

func HandleHome(context *router.Context) {
	view := view.New(context)
	view.Text("Hello World")
	view.Render(context.Writer)
}

func main() {

	// Setup server
	server, err := server.New()
	if err != nil {
		fmt.Printf("Error creating server %s", err)
		return // We may get callback to  if this is a first time setup
	}

	// Setup logger
	server.Logger = log.New(server.Config("log"), server.Production())

	// Setup our view templates
	setupView(server)

	// Setup our database
	setupDatabase(server)

	// Set up routes
	setupRouter(server)

	// Inform user of server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

	// Start the server
	err = server.Start()
	if err != nil {
		server.Fatalf("Error starting server %s", err)
	}

}

func setupRouter(server *server.Server) {

	// Routing
	router, err := router.New(server.Logger, server)
	if err != nil {
		server.Fatalf("Error creating router %s", err)
	}

	// Setup our router and handlers
	setupRoutes(router)

}

// SetupRoutes defines connections from paths to functions for this app
func setupRoutes(r *router.Router) {

	// Special route for home page
	r.Add("/", HandleHome)

}

func setupView(server *server.Server) {
	// Set up our own func map here if required
	//funcs := view.DefaultHelpers
	//view.Helpers = funcs
	view.Production = server.Production()

	err := view.LoadTemplates()

	if err != nil {
		server.Fatalf("Error reading templates %s", err)
	}

	server.Log("#info Parsed templates")
}

// Setup db - at present query pkg manages this...
func setupDatabase(server *server.Server) {
	config := server.Configuration()
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	// Ask query to open the database
	err := query.OpenDatabase(options)

	if err != nil {
		server.Fatalf("Error reading database %s", err)
	}

	server.Logf("#info Opened database at %s for user %s", server.Config("db"), server.Config("db_user"))

}
