package app

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/resource"
	"github.com/fragmenta/fragmenta-app/src/users"
)

var config = &resource.MockConfig{}

// TestRouter tests our routes are functioning correctly.
func TestRouter(t *testing.T) {

	logger := log.New(os.Stderr, "test:", log.Lshortfile)

	// Set up the router with mock config
	router, err := router.New(logger, config)
	if err != nil {
		t.Fatalf("app: error creating router %s", err)
	}

	// Set up routes
	SetupRoutes(router)

	// Setup our view templates (required for home route)
	err = view.LoadTemplatesAtPaths([]string{".."}, view.Helpers)
	if err != nil {
		t.Fatalf("app: error reading templates %s", err)
	}

	// Test serving the route / which should always exist
	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Test code on response
	if w.Code != http.StatusOK {
		t.Fatalf("app: error code on / expected:%d got:%d", http.StatusOK, w.Code)
	}

}

// TestAuth tests our authentication is functioning after setup.
func TestAuth(t *testing.T) {

	SetupAuth(config)

	user := &users.User{}

	// Test anon cannot access /users
	err := can.List(user, users.MockAnon())
	if err == nil {
		t.Fatalf("app: authentication block failed for anon")
	}

	// Test anon cannot edit admin user
	err = can.Update(users.MockAdmin(), users.MockAnon())
	if err == nil {
		t.Fatalf("app: authentication block failed for anon")
	}

	// Test admin can access /users
	err = can.List(user, users.MockAdmin())
	if err != nil {
		t.Fatalf("app: authentication failed for admin")
	}

	// Test admin can admin user
	err = can.Manage(user, users.MockAdmin())
	if err != nil {
		t.Fatalf("app: authentication failed for admin")
	}

}
