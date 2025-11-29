package main 

import (
	"exp/routers"
	"exp/websocket"
)


func main() {

    hub := websocket.NewHub()
    go hub.Run()

    router := routers.InitRouters(hub)
    router.Run(":8080")
}
