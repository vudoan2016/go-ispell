package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/ispell/input"
	"github.com/vudoan2016/ispell/output"
)

func main() {
	books := []string{"input/5000-sat-words.txt"}
	err := input.Init(books)
	if err != nil {
		panic(err)
	}

	// Initialize the router
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Next()
	})

	// Ready to serve
	router.GET("/", output.Respond)
	router.Run(":8081")
}
