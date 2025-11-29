package routers

import (
	"exp/websocket"

	"github.com/gin-gonic/gin"
)

func getCompletepath(url string) string {
	return "sploot/api/" + url
}

func InitRouters(hub *websocket.Hub) *gin.Engine {
	r := gin.Default()

	v1Location := r.Group(getCompletepath("location/ws"))
	{
		v1Location.GET("/register", func(c *gin.Context) {
			websocket.RegisterClient(hub, c)
		})
		v1Location.GET("/un-register", func(c *gin.Context) {
			websocket.UnregisterClient(hub, c)
		})
		// v1Location.GET("/ws", func(c *gin.Context) {
		// 	websocket.ServeWS(hub, c)
		// })
	}
	return r
}
