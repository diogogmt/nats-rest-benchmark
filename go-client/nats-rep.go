package main

import (
	"fmt"
	nats "github.com/nats-io/nats"
	"sync"
)

type Server struct {
	NatsConn *nats.EncodedConn
}

func (server *Server) NatsListItems(itemFilter *Item) {
	go func() {
		itemList := ItemList{
			Items: getItems(),
			Options: Options{
				RequestUUID: itemFilter.Options.RequestUUID,
			},
		}
		server.NatsConn.Publish(itemFilter.Options.Reply, itemList)
	}()
}

func (server *Server) NatsItemDetails(itemFilter *Item) {
	go func() {
		item, err := itemFilter.find()
		if err != nil {
			errResponse := Error{
				HttpCode: 400,
				Message: err.Error(),
				Options: Options{
					RequestUUID: itemFilter.Options.RequestUUID,
				},
			}
			server.NatsConn.Publish(itemFilter.Options.Reply, errResponse)
			return
		}
		item.Options = Options{
			RequestUUID: itemFilter.Options.RequestUUID,
		}
		server.NatsConn.Publish(itemFilter.Options.Reply, item)
	}()
}

func startNats() {
	// Cluster
//	var servers = "nats://nats01.cluster.local:4222, nats://nats02.cluster.local:4223, nats://nats03.cluster.local:4224"
	var servers = "nats://nats.server.local:4222"
	nc, _ := nats.Connect(servers,
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
	natsConnection, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
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
