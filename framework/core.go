package framework

import (
	"log"
	"net/http"
)

type Core struct {
	router map[string]ControllerHandler
}

func (core *Core) Get(key string, handler ControllerHandler) {
	core.router[key] = handler
}

func NewCore() *Core {
	return &Core{
		router: map[string]ControllerHandler{},
	}
}

func (core *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("web server start")
	ctx := NewContext(request, response)

	key := request.URL.Path[1:]
	handler := core.router[key]
	if handler == nil {
		return
	}
	log.Println("core.router")
	handler(ctx)
}
