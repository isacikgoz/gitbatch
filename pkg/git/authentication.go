package git

import (
	"net/url"
)

// Credentials holds user credentials to authenticate and authorize while
// communicating with remote if required
type Credentials struct {
	// User is the user id for authentication
	User string
	// Password is the secret information required for authetntication
	Password string
}

var (
	authProtocolHttp  = "http"
	authProtocolHttps = "https"
	authProtocolSSH   = "ssh"
)

// authentication protocol returns the type of protocol for given remote's URL
// various auth protocols require different kind of authentication
func authProtocol(r *Remote) (p string, err error) {
	u, err := url.Parse(r.URL[0])
	if err != nil {
		return p, err
	}
	return u.Scheme, err
}
