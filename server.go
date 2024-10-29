package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
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

	r.GET("/phase-change-diagram", func(c *gin.Context) {
		pressure := c.Query("pressure")
		if pressure == "" {
			c.JSON(400, gin.H{"status": "error", "message": "Pressure query parameter is missing"})
			return
		}
		p, err := strconv.ParseFloat(pressure, 64)
		if err != nil {
			c.JSON(400, gin.H{"status": "error", "message": "Invalid pressure value"})
			return
		}

		y1 := 0.0035
		x1 := 10.0
		y2 := 0.00105
		x2 := 0.05
		y3 := 30.0
		x3 := 0.05

		saturatedLiquidLine := ((y1-y2)/(x1-x2))*p + y1 - ((y1-y2)/(x1-x2))*x1
		saturatedVaporLine := ((y1-y3)/(x1-x3))*p + y1 - ((y1-y3)/(x1-x3))*x1

		c.JSON(200, gin.H{
			"specific_volume_liquid": roundFloat(saturatedLiquidLine, 6),
			"specific_volume_vapor":  roundFloat(saturatedVaporLine, 6),
		})
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

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
