package framework

import "net/http"

type Core struct {
}

func NewCore() *Core {
	return &Core{}
}

func (core *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {

}
