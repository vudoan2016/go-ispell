package main

import (
	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/ispell/input"
	"github.com/vudoan2016/ispell/output"
)

func main() {
	src := []input.Source{input.Source{Title: "5000 SAT words", File: "input/5000-sat-words.txt", Fn: input.ProcessVocab},
		input.Source{Title: "White Fang", File: "input/White Fang.txt", Fn: input.ProcessBook}}

	//var deck map[string]output.Vocabulary

	_, selects, err := input.Init(src)
	if err != nil {
		panic(err)
	}

	// Initialize the router
	router := gin.Default()

	router.Use(func(ctx *gin.Context) {
		ctx.Next()
	})

	// Ready to serve
	router.GET("/", output.Respond(&selects))
	router.Run(":8081")
}
