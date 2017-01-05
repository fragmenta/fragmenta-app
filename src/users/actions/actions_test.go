package useractions

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-app/src/lib/resource"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// names is used to test setting and getting the first string field of the user.
var names = []string{"foo", "bar"}

// TestSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		t.Fatalf("users: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	// Delete all users to ensure we get consistent results?
	query.ExecSQL("delete from users;")
	query.ExecSQL("ALTER SEQUENCE users_id_seq RESTART WITH 1;")
}

// Test GET /users/create
func TestShowCreateUser(t *testing.T) {

	// Create request context
	w, c := resource.GetRequestContext("/users/create", "/users/create", users.MockAdmin())

	// Run the handler
	err := HandleCreateShow(c)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("useractions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("useractions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /users/create
func TestCreateUser(t *testing.T) {
	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	// Create request context
	w, c := resource.PostRequestContext("/users/create", "/users/create", body, users.MockAdmin())

	// Run the handler to update the user
	err := HandleCreate(c)
	if err != nil {
		t.Fatalf("useractions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the user concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("useractions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the user name is in now value names[1]
	user, err := users.Find(1)
	if err != nil {
		t.Fatalf("useractions: error finding created user %s", err)
	}
	if user.ID != 1 || user.Name != names[0] {
		t.Fatalf("useractions: error with created user values: %v", user)
	}
}

// Test GET /users
func TestListUsers(t *testing.T) {

	// Create request context
	w, c := resource.GetRequestContext("/users", "/users", users.MockAdmin())

	// Run the handler
	err := HandleIndex(c)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("useractions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := "data-table-head"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("useractions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /users/1
func TestShowUser(t *testing.T) {

	// Create request context
	w, c := resource.GetRequestContext("/users/1", "/users/{id:[0-9]+}", users.MockAdmin())

	// Run the handler
	err := HandleShow(c)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("useractions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := fmt.Sprintf("<h1>%s</h1>", names[0])
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("useractions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /users/123/update
func TestShowUpdateUser(t *testing.T) {

	// Create request context
	w, c := resource.GetRequestContext("/users/1/update", "/users/{id:[0-9]+}/update", users.MockAdmin())

	// Run the handler
	err := HandleUpdateShow(c)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("useractions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("useractions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /users/123/update
func TestUpdateUser(t *testing.T) {
	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	// Create request context
	w, c := resource.PostRequestContext("/users/1/update", "/users/{id:[0-9]+}/update", body, users.MockAdmin())

	// Run the handler to update the user
	err := HandleUpdate(c)
	if err != nil {
		t.Fatalf("useractions: error handling HandleUpdateUser %s", err)
	}

	// Test we get a redirect after update (to the user concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("useractions: unexpected response code for HandleUpdateUser expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the user name is in now value names[1]
	user, err := users.Find(1)
	if err != nil {
		t.Fatalf("useractions: error finding updated user %s", err)
	}
	if user.ID != 1 || user.Name != names[1] {
		t.Fatalf("useractions: error with updated user values: %v", user)
	}

}

// Test of POST /users/123/destroy
func TestDeleteUser(t *testing.T) {

	body := strings.NewReader(``)

	// Test permissions - anon users can't destroy users

	// Create request context
	_, c := resource.PostRequestContext("/users/2/destroy", "/users/{id:[0-9]+}/destroy", body, users.MockAnon())

	// Run the handler to test failure as anon
	err := HandleDestroy(c)
	if err == nil { // failure expected
		t.Fatalf("useractions: unexpected response for HandleDestroy as anon, expected failure")
	}

	// Now test deleting the user created above as admin
	// Create request context
	w, c := resource.PostRequestContext("/users/1/destroy", "/users/{id:[0-9]+}/destroy", body, users.MockAdmin())

	// Run the handler
	err = HandleDestroy(c)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("useractions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("useractions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}

}
