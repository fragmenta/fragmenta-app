package useractions

import (
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleCreateShow serves the create form via GET for users.
func HandleCreateShow(context router.Context) error {

	user := users.New()

	// Authorise
	err := can.Create(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("user", user)
	return view.Render()
}

// HandleCreate handles the POST of the create form for users
func HandleCreate(context router.Context) error {

	user := users.New()

	// Authorise
	err := can.Create(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	userParams := user.ValidateParams(params.Map(), users.AllowedParams())

	id, err := user.Create(userParams)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to the new user
	user, err = users.Find(id)
	if err != nil {
		return router.InternalError(err)
	}

	return router.Redirect(context, user.IndexURL())
}
