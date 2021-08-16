package main

import "github.com/gin-gonic/gin"

func main(){
	r := gin.Default()

	go h.run()

	r.GET("/ws", func(c *gin.Context) {
		wsHandle(c.Writer,c.Request)
	})
	r.Run("127.0.0.1:8080")
}
