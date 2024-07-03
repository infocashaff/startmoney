package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", getFile)
	r.Run()
}

func getFile(c *gin.Context) {
	c.File("file.pdf")
}
