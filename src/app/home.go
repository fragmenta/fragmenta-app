package app

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// HandleShowHome serves the home page - in a real app this might be elsewhere
func HandleShowHome(context router.Context) {
	view := view.New(context)

	view.AddKey("title", "Hello world!")
	view.Template("app/views/home.html.got")
	view.Render(context)
}
