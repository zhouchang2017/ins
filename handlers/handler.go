package handlers

import (
	"ins/client"
	"ins/types"
)

type Handler interface {
	Read(client *client.Client, url string)
	Write(data *types.MainPageData)
}
