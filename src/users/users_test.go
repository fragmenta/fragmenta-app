// Tests for the users package
package users

import (
	"testing"

	"github.com/fragmenta/fragmenta-app/src/lib/resource"
)

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("users: Setup db failed %s", err)
	}
}

// Test Create method
func TestCreateUser(t *testing.T) {
	name := "'fué ';'\""
	userParams := map[string]string{"name": name}
	id, err := New().Create(userParams)
	if err != nil {
		t.Fatalf("users: Create user failed :%s", err)
	}

	user, err := Find(id)
	if err != nil {
		t.Fatalf("users: Create user find failed")
	}

	if user.Name != name {
		t.Fatalf("users: Create user name failed expected:%s got:%s", name, user.Name)
	}

}

// Test Index (List) method
func TestListUsers(t *testing.T) {

	// Get all users (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("users: List no user found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("users: List no users found :%s", err)
	}

}

// Test Update method
func TestUpdateUser(t *testing.T) {

	// Get the last user (created in TestCreateUser above)
	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("users: Destroy no user found :%s", err)
	}
	user := results[0]

	name := "bar"
	userParams := map[string]string{"name": name}
	err = user.Update(userParams)
	if err != nil {
		t.Fatalf("users: Update user failed :%s", err)
	}

	// Fetch the user again from db
	user, err = Find(user.ID)
	if err != nil {
		t.Fatalf("users: Update user fetch failed :%s", user.Name)
	}

	if user.Name != name {
		t.Fatalf("users: Update user failed :%s", user.Name)
	}

}

// Test Destroy method
func TestDestroyUser(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("users: Destroy no user found :%s", err)
	}
	user := results[0]
	count := len(results)

	err = user.Destroy()
	if err != nil {
		t.Fatalf("users: Destroy user failed :%s", err)
	}

	// Check new length of users returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("users: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("users: Destroy user count wrong :%d", len(results))
	}

}
