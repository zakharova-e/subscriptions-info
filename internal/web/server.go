package web

import (
	"net/http"
)

//	@title			Subscriptions service
//	@version		1.0
//	@description	Test challange
//	@host			localhost:8080
//	@BasePath		/
//	@schemes		http

func Run() {
	mux := RegisterRoutes()
	err := http.ListenAndServe(":8080", CORSMiddleware(mux))
	if err != nil {
		panic(err)
	}
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
