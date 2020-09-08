package output

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vudoan2016/ispell/input"
)

const count int = 100

type Vocabulary struct {
	Word string `json:"Word"`
	Type string `json:"Type"`
	Def  string `json:"Definition"`
}

func selectWords() []Vocabulary {
	var vocabs []Vocabulary
	var i int

	rand.Seed(time.Now().UnixNano())
	l := len(input.Deck)
	for i < count {
		v := input.Deck[rand.Intn(l)]
		vocabs = append(vocabs, Vocabulary{Word: v.Word, Type: v.Type, Def: v.Def})
		i++
	}
	return vocabs
}

// Respond processes '/' route
func Respond(ctx *gin.Context) {
	switch ctx.Request.Header.Get("Accept") {
	case "application/json":
		vocabs := selectWords()
		ctx.JSON(http.StatusOK, vocabs)
	}
}
