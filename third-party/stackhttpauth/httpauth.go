/*
Package stackhttpauth provides wrappers based on the github.com/abbot/go-http-auth package.

*/
package stackhttpauth

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"

	"github.com/abbot/go-http-auth"
	"github.com/go-on/stack"
)

type AuthenticatedRequest struct {
	auth.AuthenticatedRequest
}

func (a *AuthenticatedRequest) Swap(repl interface{}) {
	*a = *(repl.(*AuthenticatedRequest))
}

type digest struct {
	secrets func(user, realm string) string
	realm   string
}

func (d *digest) Wrap(next http.Handler) http.Handler {
	authenticator := auth.NewDigestAuthenticator(d.realm, d.secrets)
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		authR := AuthenticatedRequest{*r}
		w.(stack.Contexter).Set(&authR)
		next.ServeHTTP(w, &r.Request)
	}
	return authenticator.Wrap(fn)
}

// Digest returns a wrapper that authenticates via auth.NewDigestAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func Digest(realm string, secrets func(user, realm string) string) *digest {
	return &digest{secrets, realm}
}

type basic struct {
	secrets func(user, realm string) string
	realm   string
}

func (d *basic) Wrap(next http.Handler) http.Handler {
	authenticator := auth.NewBasicAuthenticator(d.realm, d.secrets)
	fn := func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		authR := AuthenticatedRequest{*r}
		w.(stack.Contexter).Set(&authR)
		next.ServeHTTP(w, &r.Request)
	}
	return authenticator.Wrap(fn)
}

// Basic returns a wrapper that authenticates via auth.NewBasicAuthenticator
// and saves the resulting *auth.AuthenticatedRequest in the Contexter (response writer).
func Basic(realm string, secrets func(user, realm string) string) *basic {
	return &basic{secrets, realm}
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