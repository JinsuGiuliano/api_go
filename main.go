package main

import (
	"example/api_go/controller"
	mng "example/api_go/manager"
	"example/api_go/router"
	"example/api_go/server"
	"example/api_go/store"
)

func main() {

	store := store.New()

	dbm := mng.DataManager{
		DBManager: store,
	}

	manager := mng.New(dbm)

	server := server.New()

	cont := controller.New(server, manager)

	cont.InitializeOAuthGoogle()

	router.Set(server, cont)
}
