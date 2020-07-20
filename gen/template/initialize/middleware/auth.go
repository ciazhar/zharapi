package middleware

import (
	"fmt"
	"github.com/ciazhar/zharapi/gen/template/data"
	"os"
	"path/filepath"
	"text/template"
)

func InitAuth(d data.Data, funcMap map[string]interface{}) {

	fmt.Println("init auth")

	t := template.Must(template.New("queue").Funcs(funcMap).Parse(AuthTemplate))

	if _, err := os.Stat("common/middleware"); os.IsNotExist(err) {
		newPath := filepath.Join(".", "common/middleware")
		os.MkdirAll(newPath, os.ModePerm)
	}

	f, err := os.Create("common/middleware/auth.go")
	if err != nil {
		panic(err)
	}

	if err := t.Execute(f, d); err != nil {
		panic(err)
	}
}

var AuthTemplate = `
package middleware

import (
	"fmt"
	string2 "{{.Package}}/common/string"
	error2 "{{.Package}}/common/error"
	"{{.Package}}/common/rest"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func CheckAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.String(),"/list") {
			addFilter(h, w, r, []string{"ROLE_USER"})
		} else if strings.HasPrefix(r.URL.String(),"/lust") {
			allowAll(h,w,r)
		}else {
			disAllowAll(w)
		}
	})
}

func allowAll(h http.Handler, w http.ResponseWriter, r *http.Request)  {
	h.ServeHTTP(w, r)
}

func disAllowAll(w http.ResponseWriter)  {
	_ = rest.JSON(w, 403, error2.NewF("user not authorized to access this url"))
	return
}

func addFilter(h http.Handler, w http.ResponseWriter, r *http.Request, roles []string) {

	//get authorization header
	authorizationHeader := getAuthorizationHeader(r)
	if authorizationHeader == "" {
		_ = rest.JSON(w, 401, error2.NewF("authorization header not exist"))
		return
	}

	//get token
	token := getToken(authorizationHeader)
	if token=="" {
		_ = rest.JSON(w, 401, error2.NewF("token not valid"))
		return
	}

	//get claims
	claims, err := getClaims(token)
	if err!=nil {
		_ = rest.JSON(w, 401, error2.NewF("parsing token failed, "+err.Error()))
		return
	}

	//get authorities
	authorities := getAuthorities(claims)

	//auth filter
	notContains := !checkAuthorization(authorities, roles)
	if notContains {
		_ = rest.JSON(w, 403, error2.NewF("user not authorized to access this url"))
		return
	}

	h.ServeHTTP(w, r)
}

func getAuthorizationHeader(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func getToken(authorizationHeader string) string {
	if authHeader := strings.Split(authorizationHeader, " "); len(authHeader) == 2 && authHeader[0]=="Bearer" {
		return authHeader[1]
	}else {
		return ""
	}
}

func getClaims(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		if valid := token.Valid; !valid {
			return nil, fmt.Errorf("token not valid")
		}

		return []byte("ThisIsSecretForJWTHS512SignatureAlgorithmThatMUSTHave512bitsKeySize"), nil
	}); err != nil {
		return claims, err
	}
	return claims,nil
}

func getAuthorities(claims map[string]interface{}) []interface{} {
	authorities := make([]interface{}, 0)
	if claims["authorities"] != nil {
		authorities = claims["authorities"].([]interface{})
	}
	return authorities
}

func checkAuthorization(authorities []interface{}, roles []string) bool {
	for i := range roles {
		if contain := string2.Contains(authorities, roles[i]); contain {
			return true
		}
	}
	return false
}
`
