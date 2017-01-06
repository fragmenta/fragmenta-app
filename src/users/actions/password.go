package useractions

import (
	"fmt"
	"time"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/mail"
	"github.com/fragmenta/fragmenta-app/src/users"
)

const (
	// ResetLifetime is the maximum time reset tokens are valid for
	ResetLifetime = time.Hour
)

// HandlePasswordResetShow responds to GET /users/password/reset
// by showing the password reset page.
func HandlePasswordResetShow(context router.Context) error {
	// No authorisation required, just show the view
	view := view.New(context)
	view.Template("users/views/password_reset.html.got")
	return view.Render()
}

// HandlePasswordResetSend responds to POST /users/password/reset
// by sending a password reset email.
func HandlePasswordResetSend(context router.Context) error {

	// No authorisation required

	// Find the user by email (if not found let them know)
	// Find the user by hex token in the db
	email := context.Param("email")
	user, err := users.FindFirst("email=?", email)
	if err != nil {
		return router.Redirect(context, "/users/password/reset?message=invalid_email")
	}

	// Generate a random token and url for the email
	token := auth.BytesToHex(auth.RandomToken(32))

	// Update the user record with with this token
	userParams := map[string]string{
		"password_reset_token": token,
		"password_reset_at":    query.TimeString(time.Now().UTC()),
	}
	// Direct access to the user columns, bypassing validation
	user.Update(userParams)

	// Generate the url to use in our email
	url := fmt.Sprintf("%s/users/password?token=%s", context.Config("root_url"), token)

	// Send a password reset email out to this user
	emailContext := map[string]interface{}{
		"url":  url,
		"name": user.Name,
	}
	context.Logf("#info sending reset email:%s url:%s", user.Email, url)
	e := mail.New(user.Email)
	e.Subject = "Reset Password"
	e.Template = "users/views/password_reset_mail.html.got"
	err = mail.Send(e, emailContext)
	if err != nil {
		return err
	}

	// Tell the user what we have done
	return router.Redirect(context, "/users/password/sent")
}

// HandlePasswordResetSentShow responds to GET /users/password/sent
func HandlePasswordResetSentShow(context router.Context) error {
	view := view.New(context)
	view.Template("users/views/password_sent.html.got")
	return view.Render()
}

// HandlePasswordReset responds to POST /users/password?token=DEADFISH
// by logging the user in, removing the token
// and allowing them to set their password.
func HandlePasswordReset(context router.Context) error {

	token := context.Param("token")
	if len(token) < 10 || len(token) > 64 {
		return router.InternalError(fmt.Errorf("Invalid reset token"), "Invalid Token")
	}

	// Find the user by hex token in the db
	user, err := users.FindFirst("password_reset_token=?", token)
	if err != nil {
		return router.InternalError(err)
	}

	// Make sure the reset at time is less expire time
	if time.Since(user.PasswordResetAt) > ResetLifetime {
		return router.InternalError(nil, "Token invalid", "Your password reset token has expired, please request another.")
	}

	// Remove the reset token from this user
	// using direct access, bypassing validation
	user.Update(map[string]string{"password_reset_token": ""})

	// Log in the user and store in the session
	// Now save the user details in a secure cookie, so that we remember the next request
	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(context, context.Request())
	if err != nil {
		return router.InternalError(err)
	}

	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(context.Writer())
	context.Logf("#info Login success after reset for user: %d %s", user.ID, user.Email)

	// Redirect to the user update page so that they can change their password
	return router.Redirect(context, fmt.Sprintf("/users/%d/update", user.ID))
}
