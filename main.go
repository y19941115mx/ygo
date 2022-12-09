package main

import (
	"net/http"
	"ygo/framework"
)

func main() {
	core := framework.NewCore()
	registerRoute(core)
	server := http.Server{
		Addr:    ":8888",
		Handler: core,
	}
	server.ListenAndServe()
}
