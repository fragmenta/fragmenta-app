package app

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// Default static file handler, handles assets too
func todoHandler(context router.Context) error {

	//if !strings.HasPrefix(p, "/assets/") {
	//		return router.NotFoundError(nil)
	//	}

	view := view.New(context)
	view.AddKey("title", "todos")
	return view.Render()
}
