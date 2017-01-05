package app

import (
	"github.com/fragmenta/auth"
	"github.com/fragmenta/auth/can"

	"github.com/fragmenta/fragmenta-app/src/users"
)

// SetupAuth sets up the auth pkg and authorisation for users
func SetupAuth(c Config) {

	// Set up the auth package with our secrets from config
	auth.HMACKey = auth.HexToBytes(c.Config("hmac_key"))
	auth.SecretKey = auth.HexToBytes(c.Config("secret_key"))
	auth.SessionName = c.Config("session_name")

	// Enable https cookies on production server - everyone should be on https
	if c.Production() {
		auth.SecureCookies = true
	}

	// Set up our authorisation for user roles on resources using can pkg

	// Admins are allowed to manage all resources
	can.Authorise(users.Admin, can.ManageResource, can.Anything)

	// Editors may edit their user
	can.AuthoriseOwner(users.Editor, can.UpdateResource, users.TableName)
	// ...

	// Readers may edit their user
	can.AuthoriseOwner(users.Reader, can.UpdateResource, users.TableName)

}
