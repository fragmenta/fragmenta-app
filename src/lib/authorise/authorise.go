package authorise

import (
	"github.com/fragmenta/auth"
	"github.com/fragmenta/fragmenta-app/src/users"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"
)

// Resource defines the interface for models passed to authorise.PathAndResource
type Resource interface {
	OwnedBy(int64) bool
}

// Setup authentication and authorization keys for this app
func Setup(s *server.Server) {

	// Set up our secret keys which we take from the config
	// NB these are hex strings which we convert to bytes, for ease of presentation in secrets file
	c := s.Configuration()
	auth.HMACKey = auth.HexToBytes(c["hmac_key"])
	auth.SecretKey = auth.HexToBytes(c["secret_key"])
	auth.SessionName = "fragmenta-app"

	// Enable https cookies on production server - we don't have https, so don't do this
	//	if s.Production() {
	//		auth.SecureCookies = true
	//	}

}

// CurrentUserFilter returns a filter function which sets the current user on the context
func CurrentUserFilter(c router.Context) {
	u := CurrentUser(c)
	c.Set("current_user", u)
}

// Path authorises the path for the current user
func Path(c router.Context) bool {
	return PathAndResource(c, nil)
}

// PathAndResource authorises the path and resource for the current user
// if model is nil it is ignored and permission granted
func PathAndResource(c router.Context, r Resource) bool {

	// Short circuit evaluation if this is a public path
	if publicPath(c.Path()) {
		return true
	}

	user := c.Get("current_user").(*users.User)

	switch user.Role {
	case users.RoleAdmin:
		return true
	default:
		return false
	}

	return true
}

// publicPath returns true if this path should always be allowed, regardless of user role
func publicPath(p string) bool {
	switch p {
	case "/":
		return true
	case "/users/login":
		return true
	default:
		return false
	}
}
