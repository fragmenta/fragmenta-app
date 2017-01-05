package auth

import (
	"net/http/httptest"
	"testing"

	"github.com/fragmenta/auth"
)

var (
	testKey = "12353bce2bbc4efb90eff81c29dc982de9a0176b568db18a61b4f4732cadabbc"
	set     = "foo"
)

// TestAuthenticate tests storing a value in a cookie and retreiving it again.
func TestAuthenticate(t *testing.T) {

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// Setup auth with some test values - could read these from config I guess
	auth.HMACKey = auth.HexToBytes(testKey)
	auth.SecretKey = auth.HexToBytes(testKey)
	auth.SessionName = "test_session"

	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(w, r)
	if err != nil {
		t.Fatalf("auth: failed to build session")
	}

	// Write value
	session.Set(auth.SessionUserKey, set)

	// Set the cookie on the recorder
	err = session.Save(w)
	if err != nil {
		t.Fatalf("auth: failed to save session")
	}

	session.Set(auth.SessionUserKey, "bar")

	// Try with a bogus key, should fail
	auth.SecretKey = auth.HexToBytes(testKey + "bogus")
	err = session.Load(r)
	if err == nil {
		t.Fatalf("auth: failed to detect invalid key")
	}

	// TODO: Copy the Cookie over to a new Request and build another session
	/*
		r = httptest.NewRequest("GET", "/", nil)
		r.Header.Set("Cookie", strings.Join(w.HeaderMap["Set-Cookie"], ""))
		session, err = auth.Session(w, r)
		if err != nil {
			t.Fatalf("auth: failed to build session")
		}

		got := session.Get(auth.SessionUserKey)
		if got != set {
			t.Fatalf("auth: failed to get cookie expected:%s got:%s", set, got)
		}
	*/

}
