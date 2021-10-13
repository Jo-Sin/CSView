package main

import (
	"net/http"

	"github.com/Jo-Sin/CSView/controllers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
)

func main() {
	router := httprouter.New()


	mc := controllers.GetMongoController(getMongoSession())	//mongo-db controller
	mc.InitializeDatabase()

	router.GET("/mongo-orders/:page", mc.GetOrders)
	router.GET("/count", controllers.GetCount)

	http.ListenAndServe("localhost:8080", router)
}


func getMongoSession() *mgo.Session {
	session, err := mgo.Dial("mongodb://127.0.0.1:27017")

	if err != nil {
		panic(err)
	}

	return session
}
