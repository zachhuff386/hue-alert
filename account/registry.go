package account

import (
	"reflect"
)

var (
	authTypes      = map[string]int{}
	authRegistry   = map[string]reflect.Type{}
	clientRegistry = map[string]reflect.Type{}
	handlers       = []func(){}
)

func Register(name string, typ int, auth interface{},
	client interface{}, handler func()) {

	authTypes[name] = typ
	authRegistry[name] = reflect.TypeOf(auth)
	clientRegistry[name] = reflect.TypeOf(client)
	handlers = append(handlers, handler)
}

func InitAccounts() {
	Authenticated = make(chan bool)

	for _, handler := range handlers {
		handler()
	}
}
