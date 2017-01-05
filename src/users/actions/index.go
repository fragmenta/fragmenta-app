package useractions

import (
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/fragmenta/fragmenta-app/src/lib/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
)

// HandleIndex displays a list of users.
func HandleIndex(context router.Context) error {

	// Authorise list user
	err := can.List(users.New(), auth.CurrentUser(context))
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Build a query
	q := users.Query()

	// Order by required order, or default to id asc
	switch context.Param("order") {

	case "1":
		q.Order("created desc")

	case "2":
		q.Order("updated desc")

	case "3":
		q.Order("name asc")

	default:
		q.Order("id asc")

	}

	// Filter if requested
	filter := context.Param("filter")
	if len(filter) > 0 {
		q.Where("name ILIKE ?", filter)
	}

	// Fetch the users
	results, err := users.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("filter", filter)
	view.AddKey("users", results)
	return view.Render()
}
