package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/learn-systemverilog/learn-systemverilog-api/api"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.GET("/transpile", api.Transpile)

	if err := r.Run(); err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}
