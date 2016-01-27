package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats"
	"log"
)

func startNats() {
	fmt.Println("Client waiting for connections...")
	serverUrl := "nats://nats.server.local:4222"
	fmt.Println("nats server URL: %v+", serverUrl)
	natsConnection, _ := nats.Connect(serverUrl)
	c, _ := nats.NewEncodedConn(natsConnection, nats.JSON_ENCODER)
	defer c.Close()

	natsConnection.Subscribe("listItems", func(msg *nats.Msg) {
		itemsRes := ItemListResponse{
			Items: getItems(),
		}

		res, err := json.Marshal(itemsRes)
		if err != nil {
			log.Fatalln("Error marshaling JSON items %+v", err)
		}
		natsConnection.Publish(msg.Reply, res)
	})

	c.Subscribe("itemDetails", func(msg *nats.Msg) {
		var itemReq Item
		err := json.Unmarshal(msg.Data, &itemReq)
		if err != nil {
			msgRes := fmt.Sprintf("Error decoding request %s", itemReq.Id)
			fmt.Println(msgRes)
			errStruct := ErrorResponse{Err: ErrorResponseBody{HttpCode: 400, Message: msgRes}}
			errRes, err := json.Marshal(errStruct)
			if err != nil {
				// If the json marshalling fails what res to send back to the client ?
				return
			}
			natsConnection.Publish(msg.Reply, errRes)
			return
		}

		item, err := itemReq.find()
		if err != nil {
			errStruct := ErrorResponse{Err: ErrorResponseBody{HttpCode: 400, Message: err.Error()}}
			errRes, err := json.Marshal(errStruct)
			if err != nil {
				// do something about it
			}
			natsConnection.Publish(msg.Reply, errRes)
			return
		}
		itemRes, err := json.Marshal(item)
		if err != nil {
			// do something about it
		}
		natsConnection.Publish(msg.Reply, itemRes)
	})

	donech := make(chan bool, 1)
	<-donech
}
