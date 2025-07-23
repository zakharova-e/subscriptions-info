package web

import (
	"log"
	"net/http"
	"time"
)

//	@title			Subscriptions service
//	@version		1.0
//	@description	Test challange
//	@host			localhost:8080
//	@BasePath		/
//	@schemes		http

func Run() {
	mux := RegisterRoutes()
	err := http.ListenAndServe(":8080", LogMiddleware(CORSMiddleware(mux)))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started")
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: request [%s] %s %s \n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
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
