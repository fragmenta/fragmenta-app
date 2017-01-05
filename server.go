package main

import (
	"fmt"

	"github.com/fragmenta/server"

	"github.com/fragmenta/fragmenta-app/src/app"
)

// Main entrypoint for the server which performs bootstrap, setup
// then runs the server. Most setup is delegated to the src/app pkg.
func main() {

	// Bootstrap if required (no config file found).
	if app.RequiresBootStrap() {
		err := app.Bootstrap()
		if err != nil {
			fmt.Printf("Error bootstrapping server %s\n", err)
			return
		}
	}

	// Setup our server from config
	s, err := SetupServer()
	if err != nil {
		fmt.Printf("server: error setting up %s\n", err)
		return
	}

	// Start the server
	err = s.Start()
	if err != nil {
		s.Fatalf("server: error starting %s\n", err)
	}

}

// SetupServer reads the config and sets the server up by calling app.Setup().
func SetupServer() (*server.Server, error) {

	// Setup server
	s, err := server.New()
	if err != nil {
		return nil, err
	}

	// Call the app to perform additional setup
	app.Setup(s)

	return s, nil
}
