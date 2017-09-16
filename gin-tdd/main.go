package main

import "github.com/gin-gonic/gin"

var router *gin.Engine

func main() {
	router = gin.Default()

	router.LoadHTMLGlob("templates/*")

	initializeRoutes()

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "Home Page",
	// 	})
	// })

	router.Run(":8080")
}
