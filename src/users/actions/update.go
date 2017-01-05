package useractions

import (
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleUpdateShow renders the form to update a user.
func HandleUpdateShow(context router.Context) error {

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update user
	err = can.Update(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("user", user)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a user
func HandleUpdate(context router.Context) error {

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update user
	err = can.Update(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the user from params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	userParams := user.ValidateParams(params.Map(), users.AllowedParams())

	err = user.Update(userParams)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to user
	return router.Redirect(context, user.ShowURL())
}
