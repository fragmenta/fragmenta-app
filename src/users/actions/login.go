package useractions

import (
	"fmt"

	authenticate "github.com/fragmenta/auth"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleLoginShow shows the page at /users/login
func HandleLoginShow(context router.Context) error {

	// Check they're not logged in already.
	if !auth.CurrentUser(context).Anon() {
		return router.Redirect(context, "/?warn=already_logged_in")
	}

	// Show the login page, with login failure warnings.
	view := view.New(context)
	switch context.Param("error") {
	case "failed_email":
		view.AddKey("warning", "Sorry, we couldn't find a user with that email.")
	case "failed_password":
		view.AddKey("warning", "Sorry, the password was incorrect, please try again.")
	}
	return view.Render()
}

// HandleLogin responds to POST /users/login
// by setting a cookie on the request with encrypted user data.
func HandleLogin(context router.Context) error {

	// Check they're not logged in already if so redirect.
	if !auth.CurrentUser(context).Anon() {
		return router.Redirect(context, "/?warn=already_logged_in")
	}

	// Get the user details from the database
	params, err := context.Params()
	if err != nil {
		return router.NotFoundError(err)
	}

	// Fetch the first user
	user, err := users.FindFirst("email=?", params.Get("email"))
	if err != nil {
		context.Logf("#error Login failed for user no such user : %s %s", params.Get("email"), err)
		return router.Redirect(context, "/users/login?error=failed_email")
	}

	// Check password against the stored password
	err = authenticate.CheckPassword(params.Get("password"), user.PasswordHash)
	if err != nil {
		context.Logf("#error Login failed for user : %s %s", params.Get("email"), err)
		return router.Redirect(context, "/users/login?error=failed_password")
	}

	// Now save the user details in a secure cookie, so that we remember the next request
	session, err := authenticate.Session(context, context.Request())
	if err != nil {
		context.Logf("#error problem retrieving session")
	}

	// Success, log it and set the cookie with user id
	context.Logf("#info Login success for user: %d %s", user.ID, user.Email)
	session.Set(authenticate.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(context)

	// Redirect - ideally here we'd redirect to their original request path
	return router.Redirect(context, "/")
}
