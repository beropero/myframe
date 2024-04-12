package main

import (
	"net/http"

	"kilon"
)

func main() {
	r := kilon.Default()
	r.GET("/", func(c *kilon.Context) {
		c.String(http.StatusOK, "Hello World\n")
	})
	
	r.GET("/panic", func(c *kilon.Context) {
		names := []string{"Hello_World"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8080")
}