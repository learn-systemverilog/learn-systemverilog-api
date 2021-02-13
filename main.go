package main

import (
	"github.com/gin-gonic/gin"
	"github.com/learn-systemverilog/learn-systemverilog-api/transpiler"
	"log"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	err := transpiler.Run("abacaba")
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Run(); err != nil {
		log.Fatal("Listen and serve: ", err)
	}
}
