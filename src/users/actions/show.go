package useractions

import (
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleShow displays a single user.
func HandleShow(context router.Context) error {

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise access
	err = can.Show(user, auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	w := context.Writer()
	// Set cache control headers
	w.Header().Set("Cache-Control", "no-cache, public")
	w.Header().Set("Etag", user.CacheKey())
	view := view.New(context)
	view.AddKey("user", user)
	return view.Render()
}
