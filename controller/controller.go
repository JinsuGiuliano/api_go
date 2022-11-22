package controller

import (
	store "example/api_go/manager"
	"example/api_go/server"

	"github.com/sasha-s/go-deadlock"
)

const (
	responseKey     = "responseKey"
	sessionTokenKey = "sessionTokenKey"
	qTokenKey       = "AccessToken"
)

type Controller struct {
	api       *server.Server
	data      *store.DataManager
	authMutex deadlock.Mutex
}

func New(api *server.Server, data *store.DataManager) *Controller {
	return &Controller{
		api:  api,
		data: data,
	}
}
