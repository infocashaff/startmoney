package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/file.pdf", getFile)
	r.GET("/f.png", getFPng)
	r.GET("/s.png", getSPng)
	r.Run()
}

func getFile(c *gin.Context) {
	c.File("file.pdf")
}

func getFPng(c *gin.Context) {
	c.File("f.png")
}

func getSPng(c *gin.Context) {
	c.File("s.png")
}
