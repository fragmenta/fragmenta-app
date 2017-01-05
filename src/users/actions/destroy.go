package useractions

import (
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleDestroy responds to /users/n/destroy by deleting the user.
func HandleDestroy(context router.Context) error {

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise destroy user
	err = can.Destroy(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the user
	user.Destroy()

	// Redirect to users root
	return router.Redirect(context, user.IndexURL())
}
