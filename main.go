package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/ispell/input"
	"github.com/vudoan2016/ispell/output"
)

func main() {
	//books := []string{"input/White Fang.txt"}
	vocabs := []string{"input/5000-sat-words.txt"}
	var deck []output.Vocabulary

	deck, err := input.Init(vocabs)
	if err != nil {
		panic(err)
	}

	// Initialize the router
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Next()
	})

	// Ready to serve
	router.GET("/", output.Respond(&deck))
	router.Run(":8081")
}
