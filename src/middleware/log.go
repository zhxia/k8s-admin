package middleware

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

//用于记录所有的请求日志
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Info(r.RemoteAddr, ",", r.RequestURI)
		if r.Method == "POST" {
			data, _ := ioutil.ReadAll(r.Body)
			log.Info("post data:", string(data))
			r.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		}
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}
