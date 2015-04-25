package stackhttpauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-on/stack"
	"github.com/go-on/stack/third-party/stackhttpauth"
	"gopkg.in/go-on/go-http-auth.v1"
)

func secretBasic(user, realm string) string {
	if user == "john" {
		salt := "dlPL2MqExdf"
		magic := "$1343$"
		return stackhttpauth.BasicSecret("hello", salt, magic)
	}
	return ""
}

func secretDigest(user, realm string) string {
	if user == "john" {
		return stackhttpauth.DigestSecret("john", "hello", "example.com")
	}
	return ""
}

func app(ctx stack.Contexter, rw http.ResponseWriter, req *http.Request) {
	var authReq stackhttpauth.AuthenticatedRequest
	ctx.Get(&authReq)
	rw.Write([]byte("user " + authReq.Username + " authenticated"))
}

func ExampleBasic() {

	// check that the context fulfills all requirements

	stackBasic := stack.New().
		UseWrapper(stackhttpauth.Basic("example.com", secretBasic)).
		WrapFuncWithContext(app)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.SetBasicAuth("john", "hello")
	stackBasic.ServeHTTP(rec, req)

	fmt.Println("-- success --")
	fmt.Printf("code %d\n", rec.Code)
	fmt.Println(rec.Body.String())

	rec = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/", nil)
	stackBasic.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Printf("code %d\n", rec.Code)

	// Output:
	// -- success --
	// code 200
	// user john authenticated
	// -- fail --
	// code 401
}

func ExampleDigest() {

	stackDigest := stack.New().
		UseWrapper(stackhttpauth.Digest("example.com", secretDigest)).
		WrapFuncWithContext(app)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	stackDigest.ServeHTTP(rec, req)
	fmt.Println("-- fail --")
	fmt.Printf("code %d\n", rec.Code)

	authServerHeader := digestMap(rec.Header().Get("WWW-Authenticate"))

	req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", authHeaderForClient(authServerHeader))

	rec = httptest.NewRecorder()

	stackDigest.ServeHTTP(rec, req)

	fmt.Println("-- success --")
	fmt.Printf("code %d\n", rec.Code)
	fmt.Println(rec.Body.String())

	// Output:
	// -- fail --
	// code 401
	// -- success --
	// code 200
	// user john authenticated

}

func authHeaderForClient(authServerHeader map[string]string) string {
	HA1 := auth.H("john:" + authServerHeader["realm"] + ":hello")
	HA2 := auth.H("GET:/")
	nc, cnonce := "0", "NjE4MTM2"
	response := auth.H(strings.Join([]string{HA1, authServerHeader["nonce"], nc, cnonce, authServerHeader["qop"], HA2}, ":"))

	return fmt.Sprintf(
		`Digest username="%s", realm="%s", nonce="%s", uri="%s", qop="%s", opaque="%s", response="%s", algorithm="MD5",  nc=%s, cnonce="%s"`,
		"john",
		authServerHeader["realm"],
		authServerHeader["nonce"],
		"/",
		authServerHeader["qop"],
		authServerHeader["opaque"],
		response,
		nc,
		cnonce,
	)

}

func digestMap(header string) map[string]string {
	result := map[string]string{}
	idx := strings.Index(header, "Digest")
	digeststr := strings.TrimSpace(header[idx+7:])

	for _, pair := range strings.Split(digeststr, ",") {
		a := strings.SplitN(pair, "=", 2)
		result[strings.TrimSpace(a[0])] = strings.Trim(strings.TrimSpace(a[1]), `"`)
	}
	return result
}
