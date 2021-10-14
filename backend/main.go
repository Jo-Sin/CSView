package main

import (
	"net/http"
	"github.com/Jo-Sin/CSView/backend/controllers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {
	router := httprouter.New()

	// Initialize databases and their controllers
	mc := controllers.GetMongoController(getMongoSession())
	mc.InitializeDatabase()
	getPostgresSession()
	

	// Set API paths and associated functions
	router.GET("/mongo-orders/:page", mc.GetOrders)
	router.GET("/count", controllers.GetCount)
	router.GET("/post-orders/:page", controllers.GetPostOrders)

	// Run server to listen for API requests from front-end
	http.ListenAndServe("localhost:8080", router)


	// Close DB on exit
	defer controllers.CloseDB()
}



func getMongoSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	return session
}


func getPostgresSession() {
	dbtype := "postgres"
	user := "postgres"
	password := "postgres"
	dbname := "postgres"
	host := "localhost"
	port := "5432"
	controllers.GetPostgresSession(dbtype, user, password, dbname, host, port)
}