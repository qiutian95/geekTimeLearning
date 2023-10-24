package middleware

import (
	"log"
	"mime"
	"net/http"
)

// 中间件使用方式的一种，传入一个参数，返回一个参数
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("rece a %s request from %s", r.Method, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func Validating(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		mediatype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if mediatype != "application/json" {
			http.Error(w, "invalid Contype-Type", http.StatusUnsupportedMediaType)
			return
		}
		next.ServeHTTP(w, r)
	})
}
