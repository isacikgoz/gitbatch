package git

import (
	"net/url"
	"strings"
)

// Credentials holds user credentials to authenticate and authorize while
// communicating with remote if required
type Credentials struct {
	// User is the user id for authentication
	User string
	// Password is the secret information required for authentication
	Password string
}

// Schemes for authentication
const (
	AuthProtocolHTTP  = "http"
	AuthProtocolHTTPS = "https"
	AuthProtocolSSH   = "ssh"
)

// AuthProtocol returns the type of protocol for given remote's URL
// various auth protocols require different kind of authentication
func AuthProtocol(r *Remote) (p string, err error) {
	ur := r.URL[0]
	if strings.HasPrefix(ur, "git@") {
		return "ssh", nil
	}
	u, err := url.Parse(ur)
	if err != nil {
		return p, err
	}
	return u.Scheme, err
}
