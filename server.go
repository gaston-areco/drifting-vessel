package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	r := gin.Default()

	systems := []string{"navigation", "communications", "life_support", "engines", "deflector_shield"}
	systemsCode := []string{"NAV-01", "COM-02", "LIFE-03", "ENG-04", "SHLD-05"}
	damagedSystemIndex := 0

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"damaged_system": systems[damagedSystemIndex],
		})
	})

	r.GET("/repair-bay", func(c *gin.Context) {
		htmlContent, err := os.ReadFile("repair_bay.html")
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Unable to load HTML file"})
			return
		}

		modifiedHtmlContent := strings.ReplaceAll(string(htmlContent), "{{damagedSystem}}", systemsCode[damagedSystemIndex])

		c.Data(200, "text/html; charset=utf-8", []byte(modifiedHtmlContent))
	})

	r.POST("/teapot", func(c *gin.Context) {
		c.JSON(418, gin.H{
			"message": "I'm a teapot",
		})
	})

	r.POST("/set-damaged", func(c *gin.Context) {
		var req struct {
			System string `json:"system"`
		}
		if err := c.ShouldBindJSON(&req); err == nil {
			tempIndex := containsSystem(systems, req.System)
			if tempIndex != -1 {
				damagedSystemIndex = tempIndex
				c.JSON(200, gin.H{"status": "success"})
			} else {
				c.JSON(400, gin.H{"status": "error", "message": "Invalid system"})
			}
		} else {
			c.JSON(400, gin.H{"status": "error", "message": "Invalid request"})
		}
	})

	r.Run(fmt.Sprintf(":%s", port))
}

func containsSystem(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}
