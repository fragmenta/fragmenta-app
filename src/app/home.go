package app

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// HandleShowHome serves our home page with a simple template.
// This function might be moved over to src/pages if you have a pages resource.
func homeHandler(context router.Context) error {
	view := view.New(context)
	view.AddKey("title", "Fragmenta app")
	view.Template("app/views/home.html.got")
	return view.Render()
}
