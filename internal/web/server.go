package web

import (
	"net/http"
)

func Run(){
	mux := RegisterRoutes()
	err := http.ListenAndServe(":8080",mux)
	if err!= nil{
		panic(err)
	}
}