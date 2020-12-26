package output

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	count   int    = 100
	api_key string = "Token a0ece2037f62563e2f38d2099b31fbc5624b11ab"
)

type Vocabulary struct {
	Word    string `json:"Word"`
	Type    string `json:"Type"`
	Def     string `json:"Definition"`
	Example string `json:"Example"`
}

type owlBotDefinitions struct {
	Defs []owlBotDefinition `json:"definitions"`
}

type owlBotDefinition struct {
	Definition string `json:"definition"`
	Example    string `json:"example"`
}

func getExample(word string) string {
	var example string

	// Create a new request using http
	req, err := http.NewRequest("GET", "https://owlbot.info/api/v4/dictionary/"+word, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", api_key)

	// Send req using http client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERRO] -", err)
	}

	var definitions owlBotDefinitions
	responseData, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal([]byte(responseData), &definitions)
		if err != nil {
			log.Println("Failed to unmarshal response for ", word)
		}
		if len(definitions.Defs) > 0 {
			example = definitions.Defs[0].Example
		}
	}
	return example
}

func selectWords(vocabs *[]Vocabulary) ([]Vocabulary, error) {
	var i int
	var out []Vocabulary

	rand.Seed(time.Now().UnixNano())
	l := len(*vocabs)
	for i < count {
		v := (*vocabs)[rand.Intn(l)]
		out = append(out, Vocabulary{Word: v.Word, Type: v.Type, Def: v.Def, Example: getExample(v.Word)})
		i++
	}
	return out, nil
}

// Respond processes '/' route
func Respond(vocabs *[]Vocabulary) gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		out, err := selectWords(vocabs)
		if err == nil {
			switch ctx.Request.Header.Get("Accept") {
			case "application/json":
				ctx.JSON(http.StatusOK, out)
			}
		}
	}

	return gin.HandlerFunc(handler)
}
