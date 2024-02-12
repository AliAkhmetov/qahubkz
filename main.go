package main

import (
	"log"

	"github.com/heroku/go-getting-started/repository"
	"github.com/heroku/go-getting-started/server"
	"github.com/heroku/go-getting-started/service"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	newDbName       = "./st.db"
	initSqlFileName = "./init-up.sql"
)

func main() {

	// New store instance
	storage, err := repository.New(newDbName)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	// Init DB by init-up.sql
	if err := storage.Init(initSqlFileName); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	//New Repository struct with interfaces
	repos := repository.NewRepository(storage.Db)

	//New Handler struct
	handler := server.NewHandler(repos)
	//Create Super User
	service.CreateSuperUser(repos)
	server.Server(handler)
}
