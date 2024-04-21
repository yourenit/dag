package main

import (
	"io"
	"log"

	"github.com/yourenit/dag_engine/engine"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/simulate", func(c *gin.Context) {
		e := engine.Engine{}
		body,err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Fatal(err)
		}

		res,err := e.Evaluate(string(body))
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, res)
	})
	// r.Use()
	r.Static("/", "./static")
	r.Run(":3000")
}