package app

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
)

// Serve static files (assets, images etc)
func fileHandler(context router.Context) error {

	// First try serving assets
	err := serveAsset(context)
	if err == nil {
		return nil
	}

	// If assets fail, try to serve file in public
	return serveFile(context)
}

// serveFile serves a file from ./public if it exists
func serveFile(context router.Context) error {

	// Try a local path in the public directory
	localPath := "./public" + path.Clean(context.Path())
	s, err := os.Stat(localPath)
	if err != nil {
		// If file not found return 404
		if os.IsNotExist(err) {
			return router.NotFoundError(err)
		}

		// For other errors return not authorised
		return router.NotAuthorizedError(err)
	}

	// If not a file return immediately
	if s.IsDir() {
		return nil
	}

	// If the file exists and we can access it, serve it with cache control
	context.Writer().Header().Set("Cache-Control", "max-age:3456000, public")
	http.ServeFile(context, context.Request(), localPath)
	return nil
}

// serveAsset serves a file from ./public/assets usings appAssets
func serveAsset(context router.Context) error {
	p := path.Clean(context.Path())

	// It must be under /assets, or we don't serve
	if !strings.HasPrefix(p, "/assets/") {
		return router.NotFoundError(nil)
	}

	// Try to find an asset in our list
	f := appAssets.File(path.Base(p))
	if f == nil {
		return router.NotFoundError(nil)
	}

	// Serve the local file, with cache control
	localPath := "./" + f.LocalPath()
	context.Writer().Header().Set("Cache-Control", "max-age:3456000, public")
	http.ServeFile(context, context.Request(), localPath)
	return nil
}

// errHandler renders an error using error templates if available
func errHandler(context router.Context, e error) {

	// Cast the error to a status error if it is one, if not wrap it in a Status 500 error
	err := router.ToStatusError(e)
	context.Logf("#error %s\n", err)

	view := view.New(context)
	view.AddKey("title", err.Title)
	view.AddKey("message", err.Message)
	// In production, provide no detail for security reasons
	if !context.Production() {
		view.AddKey("status", err.Status)
		view.AddKey("file", err.FileLine())
		view.AddKey("error", err.Err)
	}
	view.Template("app/views/error.html.got")
	context.Writer().WriteHeader(err.Status)
	view.Render()
}
