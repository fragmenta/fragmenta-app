// Package users represents the user resource
package users

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/fragmenta/fragmenta-app/src/lib/resource"
	"github.com/fragmenta/fragmenta-app/src/lib/status"
)

// User handles saving and retreiving users from the database
type User struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base
	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	Email             string
	EncryptedPassword string
	Name              string
	Role              int64
}

const (
	// TableName is the database table for this resource
	TableName = "users"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"status", "email", "name", "role"}
}

// NewWithColumns creates a new user instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *User {

	user := New()
	user.ID = resource.ValidateInt(cols["id"])
	user.CreatedAt = resource.ValidateTime(cols["created_at"])
	user.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	user.Status = resource.ValidateInt(cols["status"])
	user.Email = resource.ValidateString(cols["email"])
	user.EncryptedPassword = resource.ValidateString(cols["encrypted_password"])
	user.Name = resource.ValidateString(cols["name"])
	user.Role = resource.ValidateInt(cols["role"])

	return user
}

// New creates and initialises a new user instance.
func New() *User {
	user := &User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.TableName = TableName
	user.KeyName = KeyName
	user.Status = status.Draft
	return user
}

// Find fetches a single user record from the database by id.
func Find(id int64) (*User, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all user records matching this query from the database.
func FindAll(q *query.Query) ([]*User, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of users constructed from the results
	var users []*User
	for _, cols := range results {
		p := NewWithColumns(cols)
		users = append(users, p)
	}

	return users, nil
}
