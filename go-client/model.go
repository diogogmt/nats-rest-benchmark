package main

import (
	"errors"
)

type ErrorResponse struct {
	Err ErrorResponseBody `json:"err"`
}

type ErrorResponseBody struct {
	Message  string `json:"msg"`
	HttpCode int    `json:"httpCode"`
}

type ItemListResponse struct {
	Items []Item `json:"items"`
}

type Item struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Brand string `json:"brand"`
}

func (item *Item) find() (Item, error) {
	items := getItems()
	for i := 0; i < len(items); i++ {
		if items[i].Id == item.Id {
			return items[i], nil
		}
	}
	return Item{}, errors.New("Item with ID " + item.Id + " not found")
}

func getItems() []Item {
	item1 := Item{
		Id:    "93416ae8-7225-4bb6-bc93-31dd30dc55a6",
		Name:  "MacBook Pro",
		Brand: "Apple",
	}

	item2 := Item{
		Id:    "3637fc25-7a10-48a6-83f7-d4e432a24f44",
		Name:  "Alien Ware",
		Brand: "Dell",
	}
	return []Item{item1, item2}
}
