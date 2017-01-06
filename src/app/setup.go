package app

import (
	"time"

	"github.com/fragmenta/assets"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/mail"
	"github.com/fragmenta/fragmenta-app/src/lib/mail/adapters/sendgrid"
)

// Config is used to pass settings to setup functions.
type Config interface {
	Production() bool
	Configuration() map[string]string
	Config(string) string
}

// appAssets holds a reference to our assets for use in asset setup used in the handlers.
var appAssets *assets.Collection

// Setup sets up our application.
func Setup(server *server.Server) {

	// Setup log
	server.Logger = log.New(server.Config("log"), server.Production())

	// Set up our mail adapter
	SetupMail(server)

	// Set up our assets
	SetupAssets(server)

	// Setup our view templates
	SetupView(server)

	// Setup our database
	SetupDatabase(server)

	// Set up auth pkg and authorisation for access
	SetupAuth(server)

	// Create a new router
	router, err := router.New(server.Logger, server)
	if err != nil {
		server.Fatalf("Error creating router %s", err)
	}

	// Setup our router and handlers
	SetupRoutes(router)

	// Inform user of imminent server setup
	server.Logf("#info Starting server in %s mode on port %d", server.Mode(), server.Port())

}

// SetupMail sets us up to send mail via sendgrid (requires key).
func SetupMail(server *server.Server) {
	mail.Production = server.Production()
	mail.Service = sendgrid.New(server.Config("mail_from"), server.Config("mail_secret"))
}

// SetupAssets compiles or copies our assets from src into the public assets folder.
func SetupAssets(server *server.Server) {
	defer server.Timef("#info Finished loading assets in %s", time.Now())

	// Compilation of assets is done on deploy
	// We just load them here
	assetsCompiled := server.ConfigBool("assets_compiled")
	appAssets = assets.New(assetsCompiled)

	// Load asset details from json file on each run
	err := appAssets.Load()
	if err != nil {
		// Compile assets for the first time
		server.Logf("#info Compiling assets")
		err := appAssets.Compile("src", "public")
		if err != nil {
			server.Fatalf("#error compiling assets %s", err)
		}
	}

	// Set up helpers which are aware of fingerprinted assets
	// These behave differently depending on the compile flag above
	// when compile is set to no, they use precompiled assets
	// otherwise they serve all files in a group separately
	view.Helpers["style"] = appAssets.StyleLink
	view.Helpers["script"] = appAssets.ScriptLink

}

// SetupView sets up the view package by loadind templates.
func SetupView(server *server.Server) {
	defer server.Timef("#info Finished loading templates in %s", time.Now())

	view.Production = server.Production()
	err := view.LoadTemplates()
	if err != nil {
		server.Fatalf("Error reading templates %s", err)
	}

}

// SetupDatabase sets up the db with query given our server config.
func SetupDatabase(server *server.Server) {
	defer server.Timef("#info Finished opening in %s database %s for user %s", time.Now(), server.Config("db"), server.Config("db_user"))

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

}
