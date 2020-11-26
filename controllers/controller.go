package controller

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/DLTcollab/vehicle-data-explorer/models/device"
	"github.com/DLTcollab/vehicle-data-explorer/models/elasticsearch"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_CBCDecrypter"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_deserializer"
	"github.com/DLTcollab/vehicle-data-explorer/models/jwt"
	"github.com/DLTcollab/vehicle-data-explorer/models/obd"

	"github.com/gin-gonic/gin"
)

type MAM_subscribe struct {
	Chid        string `json:"chid"`
	Network     string `json:"network"`
	Private_key string `json:"private-key"`
}

type MAM_response struct {
	Payload []string `json:"payload"`
	Chid1   string   `json:"chid"`
}

type Endpoint_obd2_data struct {
	Timestamp int64         `json:"timestamp"`
	Device_id string        `json"device_id"`
	Data      obd.ODB2_data `json:"data"`
}

type MAM_post_data struct {
	Data_id  MAM_post_innder `json:"data_id"`
	Protocol string          `json:"protocol"`
}

type MAM_post_innder struct {
	Chid string `json:"chid"`
}

func MAM_recv(host string, chid string) *MAM_response {
	var data = MAM_post_data{
		Data_id: MAM_post_innder{
			Chid: chid,
		},
		Protocol: "MAM_V1",
	}

	b, _ := json.Marshal(data)

	resp, err := http.Post(host+"/mam/recv", "application/json", bytes.NewBuffer(b))

	if err != nil {
		log.Println(err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-OK HTTP status: ", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("Error message: ", string(body))
		return nil
	}

	var mam_response MAM_response
	json.NewDecoder(resp.Body).Decode(&mam_response)

	log.Printf("Payload length: %d, chid: %s", len(mam_response.Payload), mam_response.Chid1)

	return &mam_response
}

func MAM_sub(c *gin.Context) {
	var json_data MAM_subscribe
	if err := c.BindJSON(&json_data); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "posted",
			"message": "failed",
		})
		return
	}

	mam_response := MAM_recv(json_data.Network, json_data.Chid)

	if mam_response == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "request not found",
		})
	}

	c.SetCookie("ta_host", json_data.Network, 0, "/", "localhost", false, true)
	c.SetCookie("ch_id", json_data.Chid, 0, "/", "localhost", false, true)
	c.SetCookie("private_key", json_data.Private_key, 0, "/", "localhost", false, true)

	// response to browser
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mam_response.Payload,
	})

}

func Get_dashboard_realtime_data(c *gin.Context) {

	ta_host, ta_host_err := c.Cookie("ta_host")
	chid, chid_err := c.Cookie("ch_id")
	if ta_host_err != nil {
		ta_host = "http://localhost"
	}

	if chid_err != nil {
		chid = "DEFAULTCHANNELID"
	}

	mam_response := MAM_recv(ta_host, chid)

	if mam_response == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "request not found",
		})
	}

	var data []Endpoint_obd2_data
	private_key, _ := c.Cookie("private_key")

	for i := 0; i < len(mam_response.Payload); i++ {
		payload := mam_response.Payload[i]
		endpoint_data := Descrypt_mam_response(payload, private_key)
		data = append(data, endpoint_data)
	}

	// response to browser
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func Descrypt_mam_response(mam_message string, private_key string) Endpoint_obd2_data {

	endpoint_serial := endpoint_deserializer.Endpoint_deserializer(mam_message)
	ciphertext, _ := hex.DecodeString(endpoint_serial.Ciphertext)
	plaintext := endpoint_CBCDecrypter.Endpoint_CBCDecrypter(string(ciphertext), private_key, endpoint_serial.IV, endpoint_serial.Timestamp)

	var endpoint_data Endpoint_obd2_data
	err := json.Unmarshal([]byte(plaintext), &endpoint_data)
	if err != nil {
		fmt.Println("error:", err)
	}

	log.Println("Descrypt mam response successfully")
	return endpoint_data
}

func Register_device(c *gin.Context) {
	var regDevice device.Device
	if err := c.BindJSON(&regDevice); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Failed to parse json",
		})
		return
	}

	// concate with sha256
	hexStr := device.GetDeviceHash(regDevice)
	db := c.Keys["defaultKVDatabase"].(*DefaultKVDatabase)

	exists, _ := db.Has(hexStr)

	if exists {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Device has been registered",
		})
		return
	}

	// Set device
	JSONDevice := device.EncodeDevice(regDevice)
	db.Set(hexStr, JSONDevice)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Success to register new device",
	})
}

func GrantAccessToken(c *gin.Context) {
	type Hash struct {
		Hash string `json:"hash"`
	}

	var hash Hash
	if err := c.BindJSON(&hash); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Failed to parse json",
		})
		return
	}

	db := c.Keys["defaultKVDatabase"].(*DefaultKVDatabase)
	exists, err := db.Has(hash.Hash)
	if !exists {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Not a valid device",
		})
		return
	}

	token, err := jwt.CreateJwtToken(hash.Hash)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Success to grant access token",
		"token":   token,
	})
}

func AuthRequired(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	tokenString := strings.Split(auth, "Bearer ")[1]

	token, err := jwt.VerifyJwtToken(tokenString)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Not a valid token",
		})
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(*jwt.Claims); ok && token.Valid {
		c.Set("deviceHash", claims.Hash)
		c.Next()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "failed",
		"message": "Not a valid token",
	})
	c.Abort()
}

func InsertDeviceLog(c *gin.Context) {
	deviceHash := c.MustGet("deviceHash").(string)

	type JSONReq struct {
		Log string `json:"log"`
	}

	var req JSONReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Failed to parse json",
		})
		return
	}

	if deviceHash == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Device not found",
		})
		return
	}

	db := c.Keys["defaultKVDatabase"].(*DefaultKVDatabase)
	devStr, err := db.Get(deviceHash)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Not a valid token",
		})
		return
	}

	dev, err := device.DecodeDevice(devStr)

	if err != nil {
		log.Print(err)
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "Device not found",
		})
		return
	}

	status := elasticsearch.InsertDeviceLog(dev.IMEI, req.Log)

	if status {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Success to insert device log",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "failed",
		"message": "Failed to to Insert device log",
	})
}

func QueryDeviceLog(c *gin.Context) {

	deviceID := c.Query("deviceID")
	logs := elasticsearch.QueryDeviceLog(deviceID)
	if logs == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  "failed",
			"message": "No log found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "success to query device logs",
		"log":     logs,
	})
}
