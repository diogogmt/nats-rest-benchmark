package main

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats"
	"sync"
	// "time"
)

type Server struct {
	NatsConn *nats.Conn
}

func (server *Server) NatsListItems(msg *nats.Msg) {
	// fmt.Printf("NATS - Processing listItems request - %s\n", time.Now())
	go func() {
		// fmt.Printf("NATS - Executing listItems request - %s\n", time.Now())

		var serviceReq ServiceRequest
		json.Unmarshal(msg.Data, &serviceReq)

		itemsRes := ItemListResponse{
			Items: getItems(),
		}
		itemsJsonRes, _ := json.Marshal(itemsRes)

		serviceRes := ServiceResponse{
			RequestUUID: serviceReq.RequestUUID,
			Data:        string(itemsJsonRes),
		}
		serviceJsonRes, _ := json.Marshal(serviceRes)
		server.NatsConn.Publish(serviceReq.Reply, serviceJsonRes)
	}()
}

func (server *Server) NatsItemDetails(msg *nats.Msg) {
	// fmt.Printf("NATS - Processing itemDetails request - %s\n", time.Now())
	go func() {
		// fmt.Printf("NATS - Executing itemDetails request - %s\n", time.Now())
		var serviceReq ServiceRequest
		json.Unmarshal(msg.Data, &serviceReq)
		dataByte := []byte(serviceReq.Data)
		var itemReq Item
		json.Unmarshal(dataByte, &itemReq)

		item, err := itemReq.find()
		if err != nil {
			errStruct := ErrorResponse{Err: ErrorResponseBody{HttpCode: 400, Message: err.Error()}}
			errJson, _ := json.Marshal(errStruct)
			serviceRes := ServiceResponse{
				RequestUUID: serviceReq.RequestUUID,
				Data:        string(errJson),
			}
			serviceResJson, _ := json.Marshal(serviceRes)
			server.NatsConn.Publish(serviceReq.Reply, serviceResJson)
			return
		}
		itemJson, _ := json.Marshal(item)

		resJson := ServiceResponse{
			RequestUUID: serviceReq.RequestUUID,
			Data:        string(itemJson),
		}
		serviceJsonRes, _ := json.Marshal(resJson)
		server.NatsConn.Publish(serviceReq.Reply, serviceJsonRes)
	}()
}

func startNats() {
	var servers = "nats://nats01.server.local:4222, nats://nats02.server.local:4223, nats://nats03.server.local:4224"
	natsConnection, _ := nats.Connect(servers,
		nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
			fmt.Printf("Nats ErrorHandler\n")
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Printf("Nats CloseHandler\n")
		}),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Nats DisconnectHandler\n")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Printf("Nats ReconnectHandler\n", nc.ConnectedUrl())
		}))
	defer natsConnection.Close()
	fmt.Println("Client waiting for connections...")

	server := Server{
		NatsConn: natsConnection,
	}

	natsConnection.Subscribe("listItems", server.NatsListItems)
	natsConnection.Subscribe("itemDetails", server.NatsItemDetails)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
