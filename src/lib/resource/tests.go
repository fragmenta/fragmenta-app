package resource

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// This file contains some test helpers for resources.

// basePath returns the path to the fragmenta root from a given test folder.
func basePath(depth int) string {
	// Construct a path to root
	p := ""
	for i := 0; i < depth; i++ {
		p = filepath.Join(p, "..")
	}
	return p
}

// SetupAuthorisation sets up mock authorisation.
func SetupAuthorisation() {
	// Set up some simple permissions for testing -
	//  at present we just test on admins if testing other permissions
	// they'd need to be added here
	can.Authorise(100, can.ManageResource, can.Anything)
}

// SetupView sets up the view package for testing by loading templates.
func SetupView(depth int) error {
	view.Production = false
	return view.LoadTemplatesAtPaths([]string{filepath.Join(basePath(depth), "src")}, view.Helpers)
}

// SetupTestDatabase sets up the database for all tests from the test config.
func SetupTestDatabase(depth int) error {

	// Read config json
	path := filepath.Join(basePath(depth), "secrets", "fragmenta.json")
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	config := data["test"]
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
	}

	// Ask query to open the database
	err = query.OpenDatabase(options)
	if err != nil {
		return err
	}

	// For speed
	query.Exec("set synchronous_commit=off;")
	return nil
}

// MockConfig conforms to the config interface.
type MockConfig struct {
	Data map[string]string
}

// Production returns false.
func (c *MockConfig) Production() bool {
	return false
}

// Config returns the config value for key.
func (c *MockConfig) Config(key string) string {
	return c.Data[key]
}

// Configuration returns the current config
func (c *MockConfig) Configuration() map[string]string {
	return c.Data
}

// TestContextForRequest returns a context for testing handlers.
func TestContextForRequest(w http.ResponseWriter, r *http.Request, pattern string) router.Context {
	route, err := router.NewRoute(pattern, nil)
	if err != nil {
		return nil
	}
	return router.NewContext(w, r, route, &MockConfig{}, log.New(os.Stderr, "test:", log.Lshortfile))
}

// GetRequestContext returns a context for testing GET handlers with a path and current user.
func GetRequestContext(path string, pattern string, u interface{}) (*httptest.ResponseRecorder, router.Context) {
	r := httptest.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	c := TestContextForRequest(w, r, pattern)
	c.Set("current_user", u) // Set user for testing permissions
	return w, c
}

// PostRequestContext returns a context for testing POST handlers with a path, body and current user.
func PostRequestContext(path string, pattern string, body io.Reader, u interface{}) (*httptest.ResponseRecorder, router.Context) {
	r := httptest.NewRequest("POST", path, body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	c := TestContextForRequest(w, r, pattern)
	c.Set("current_user", u) // Set user for testing permissions
	return w, c
}
