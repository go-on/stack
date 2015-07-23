/*
Package stackhttpauth provides wrappers based on the github.com/abbot/go-http-auth package.

*/
package stackhttpauth

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/go-on/stack.v3"
	"gopkg.in/go-on/go-http-auth.v1"
)

type AuthenticatedRequest struct {
	auth.AuthenticatedRequest
}

func (a *AuthenticatedRequest) Recover(repl interface{}) {
	*a = *(repl.(*AuthenticatedRequest))
}

type digest struct {
	auth *auth.DigestAuth
}

func (d *digest) ServeHTTP(ctx stack.Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler) {
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		authR := AuthenticatedRequest{*r}
		ctx.Set(&authR)
		next.ServeHTTP(w, req)
	}
	d.auth.Wrap(fn).ServeHTTP(wr, req)
}

// NewDigest returns a middleware that authenticates via auth.NewDigestAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func NewDigest(realm string, secrets func(user, realm string) string) *digest {
	return &digest{auth.NewDigestAuthenticator(realm, secrets)}
}

type basic struct {
	auth *auth.BasicAuth
}

func (d *basic) ServeHTTP(ctx stack.Contexter, wr http.ResponseWriter, req *http.Request, next http.Handler) {
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		authR := AuthenticatedRequest{*r}
		ctx.Set(&authR)
		next.ServeHTTP(w, req)
	}
	d.auth.Wrap(fn).ServeHTTP(wr, req)
}

// NewBasic returns a middleware that authenticates via auth.NewBasicAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func NewBasic(realm string, secrets func(user, realm string) string) *basic {
	return &basic{auth.NewBasicAuthenticator(realm, secrets)}
}

func DigestSecret(user, password, realm string) string {
	m := md5.New()
	io.WriteString(m, user+":"+realm+":"+password)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func BasicSecret(password, salt, magic string) string {
	return string(auth.MD5Crypt([]byte(password), []byte(salt), []byte(magic)))
	// md5e := auth.NewMD5Entry("$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1")
}
