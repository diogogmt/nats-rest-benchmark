package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func ListItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	itemsRes := ItemList{
		Items: getItems(),
	}
	res, err := json.Marshal(itemsRes)
	if err != nil {
		log.Fatalln("Error marshaling JSON items %+v", err)
	}
	w.Write(res)
}

func ItemDetails(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	itemReq := Item{Id: ps.ByName("id")}
	item, err := itemReq.find()
	if err != nil {
		errStruct := Error{HttpCode: 400, Message: err.Error()}
		errRes, err := json.Marshal(errStruct)
		if err != nil {
			// do something about it
			log.Fatalln("Error marshaling JSON items %+v", err)
		}
		w.Write(errRes)
		return
	}
	itemJson, err := json.Marshal(item)
	if err != nil {
		// do something about it
		log.Fatalln("Error marshaling JSON items %+v", err)
	}
	w.Write(itemJson)
}

func startRest() {
	router := httprouter.New()
	router.GET("/items", ListItems)
	router.GET("/items/:id", ItemDetails)
	fmt.Println("Waiting for HTTP connections on port 8080")
	log.Fatal(http.ListenAndServe(":4000", router))
}
