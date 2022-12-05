package main

import (
	"net/http"
	"ygo/framework"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: framework.NewCore(),
	}
	server.ListenAndServe()

}
