package output

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Vocabulary struct {
	Word  string `json:"Word"`
	Type  string `json:"Type"`
	Def   string `json:"Definition"`
	Usage string `json:"Usage"`
}

// Respond processes '/' route
func Respond(selects *[]Vocabulary) gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		switch ctx.Request.Header.Get("Accept") {
		case "application/json":
			ctx.JSON(http.StatusOK, *selects)
		}
	}

	return gin.HandlerFunc(handler)
}
