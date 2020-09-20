package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func showIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"status": "success",
	})
}

func show_dash_board(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"status": "success",
	})
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "public")
	r.Static("/assets", "public/assets")

	r.GET("/", showIndex)
	r.GET("/dashboard", show_dash_board)
	api := r.Group("/api")
	{
		api.POST("/sub_mam", MAM_sub)
		api.GET("/dashboard_data", Get_dashboard_realtime_data)
	}
	return r
}
