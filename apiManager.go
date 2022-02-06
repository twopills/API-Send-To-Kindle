package main

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()
	Post(r)
	r.Run(":3001")
}

func Post(r *gin.Engine) {
	r.POST("/kindle", func(c *gin.Context) {
		link := c.PostForm("link")
		title := c.PostForm("title")
		c.JSON(200, gin.H{
			"link":  link,
			"title": title,
		})
		if strings.Contains(link, "https") || strings.Contains(link, "http") {
			log.Println("Eccomi", link, title)
			takeOneHtmlElement(link, title)
		}
	})
}
