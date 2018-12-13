package git

import (
	"net/url"
)

// Credentials holds user credentials to authenticate and authorize while
// communicating with remote if required
type Credentials struct {
	User     string
	Password string
}

var (
	authProtocolHttp  = "http"
	authProtocolHttps = "https"
	authProtocolSSH   = "ssh"
)

func (entity *RepoEntity) authProtocol(remote *Remote) (p string, err error) {
	u, err := url.Parse(remote.URL[0])
	if err != nil {
		return p, err
	}
	return u.Scheme, err
}
