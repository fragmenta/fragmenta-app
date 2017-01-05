package users

// This file contains functions related to authorisation and roles.

// User roles
const (
	Anon   = 0
	Editor = 10
	Reader = 20
	Admin  = 100
)

// Anon returns true if this user is not a logged in user.
func (u *User) Anon() bool {
	return u.Role == Anon || u.ID == 0
}

// Admin returns true if this user is an Admin.
func (u *User) Admin() bool {
	return u.Role == Admin
}

// Reader returns true if this user is an Reader.
func (u *User) Reader() bool {
	return u.Role == Reader
}

// can.User interface

// RoleID returns the user role for auth.
func (u *User) RoleID() int64 {
	return u.Role
}

// UserID returns the user id for auth.
func (u *User) UserID() int64 {
	return u.ID
}
