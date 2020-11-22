package controller

import (
	"net/http"

	"github.com/DLTcollab/vehicle-data-explorer/models/kvstore"
	"github.com/DLTcollab/vehicle-data-explorer/models/redis"
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

func ShowRegisterDevice(c *gin.Context) {
	c.HTML(http.StatusOK, "register_device.html", gin.H{
		"status": "success",
	})
}

func ShowQueryDeviceLog(c *gin.Context) {
	c.HTML(http.StatusOK, "show_query_device_log.html", gin.H{
		"status": "success",
	})
}

type DefaultKVDatabase struct {
	kvstore.KVStore
}

func SetKVDatabase() gin.HandlerFunc {
	dbInstance := redis.New()
	db := &DefaultKVDatabase{dbInstance}

	return func(c *gin.Context) {
		c.Set("defaultKVDatabase", db)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/public", "public")
	r.Static("/assets", "public/assets")

	r.Use(SetKVDatabase())

	r.GET("/", showIndex)
	r.GET("/dashboard", show_dash_board)
	r.GET("/register_device", ShowRegisterDevice)
	r.GET("/query_device_log", ShowQueryDeviceLog)

	api := r.Group("/api")
	{
		api.POST("/sub_mam", MAM_sub)
		api.GET("/dashboard_data", Get_dashboard_realtime_data)
		api.POST("/register_device", Register_device)
		api.POST("/grant_access_token", GrantAccessToken)
		api.GET("/log/query", QueryDeviceLog)
	}

	authorized := r.Group("/api")
	authorized.Use(AuthRequired)
	{
		authorized.POST("/log/insert", InsertDeviceLog)
	}
	return r
}
