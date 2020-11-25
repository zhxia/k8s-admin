package middleware

import (
	"encoding/base64"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	AuthUsers map[string]string
}

func (am *AuthMiddleware) Parse(authMap map[string]gjson.Result) {
	if len(authMap) > 0 {
		am.AuthUsers = make(map[string]string)
		for k, v := range authMap {
			am.AuthUsers[k] = v.String()
		}
	}
	log.Debug("authentication users:", am.AuthUsers)
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println("check ...")
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			send401Status(w)
			return
		}
		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			send401Status(w)
			return
		}
		authPair := strings.SplitN(string(b), ":", 2)
		if len(authPair) != 2 {
			send401Status(w)
			return
		}
		hasAuth := false
		if password, ok := am.AuthUsers[authPair[0]]; ok {
			if password == authPair[1] {
				hasAuth = true
			}
		}
		if !hasAuth {
			send401Status(w)
			return
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func send401Status(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="MY REALM"`)
	w.WriteHeader(401)
	w.Write([]byte("401 Unauthorized\n"))
}
