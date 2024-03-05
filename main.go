package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/heroku/go-getting-started/repository"
	"github.com/heroku/go-getting-started/server"
	"github.com/heroku/go-getting-started/service"
	_ "github.com/heroku/x/hmetrics/onload"
)

const (
	newDbName       = "./st.db"
	initSqlFileName = "./init-up.sql"
	host            = "localhost"
	port            = "5432"
	user            = "postgres"
	password        = "123456"
	dbname          = "qahub"
)

func main() {

	// New store instance
	storage, err := repository.New(host, port, user, password, dbname)
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
	srv := new(server.Server)
	port := os.Getenv("PORT")

	if err := srv.Run(port, handler.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
	//server.Server(handler)
}
