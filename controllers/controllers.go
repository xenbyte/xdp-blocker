package controllers

import "github.com/gin-gonic/gin"

func Welcome(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "Welcome to XDP-Dropper Firewall!",
	})
}
