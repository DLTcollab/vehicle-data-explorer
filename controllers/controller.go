package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/DLTcollab/vehicle-data-explorer/OBD"
	"github.com/DLTcollab/vehicle-data-explorer/models/device"
	"github.com/DLTcollab/vehicle-data-explorer/models/elasticsearch"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_CBCDecrypter"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_deserializer"
	"github.com/DLTcollab/vehicle-data-explorer/models/jwt"
	"github.com/google/brotli/go/cbrotli"

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
	Data      OBD.OBD2_json `json:"data"`
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
		obd2Meta := Descrypt_mam_response(payload, private_key)
		obd2Data := new(OBD.OBD2_data)
		obd2Meta.Data(obd2Data)

		endpointData := Endpoint_obd2_data{
			Timestamp: obd2Meta.Timestamp(),
			Device_id: string(obd2Meta.DeviceID()),
			Data: OBD.OBD2_json{
				Vin:                                 string(obd2Data.Vin()),
				Engine_load:                         OBD.GetCalculatedEngineLoad(obd2Data.EngineLoad()),
				Engine_coolant_temperature:          OBD.GetEngineCoolantTemperature(obd2Data.EngineCoolantTemperature()),
				Fuel_pressure:                       OBD.GetFuelPressure(obd2Data.FuelPressure()),
				Engine_speed:                        OBD.GetEngineSpeed(obd2Data.EngineSpeed()),
				Vehicle_speed:                       OBD.GetVehicleSpeed(obd2Data.VehicleSpeed()),
				Intake_air_temperature:              OBD.GetIntakeAirTemperature(obd2Data.IntakeAirTemperature()),
				Mass_air_flow:                       OBD.GetMassAirFlow(obd2Data.MassAirFlow()),
				Fuel_tank_level_input:               OBD.GetFuelTankLevelInput(obd2Data.FuelTankLevelInput()),
				Absolute_barometric_pressure:        OBD.GetAbsoluteBarometricPressure(obd2Data.AbsoluteBarometricPressure()),
				Control_module_voltage:              OBD.GetControlModuleVoltage(obd2Data.ControlModuleVoltage()),
				Throttle_position:                   OBD.GetThrottlePosition(obd2Data.ThrottlePosition()),
				Ambient_air_temperature:             OBD.GetAmbientAirTemperature(obd2Data.AmbientAirTemperature()),
				Relative_accelerator_pedal_position: OBD.GetRelativeAcceleratorPedalPosition(obd2Data.RelativeAcceleratorPedalPosition()),
				Engine_oil_temperature:              OBD.GetEngineOilTemperature(obd2Data.EngineOilTemperature()),
				Engine_fuel_rate:                    OBD.GetEngineFuelRate(obd2Data.EngineFuelRate()),
				Service_distance:                    int(obd2Data.ServiceDistance()),
				Anti_lock_barking_active:            int(obd2Data.AntiLockBarkingActive()),
				Steering_wheel_angle:                int(obd2Data.SteeringWheelAngle()),
				Position_of_doors:                   int(obd2Data.PositionOfDoors()),
				Right_left_turn_signal_light:        int(obd2Data.RightLeftTurnSignalLight()),
				Alternate_beam_head_light:           int(obd2Data.AlternateBeamHeadLight()),
				High_beam_head_light:                int(obd2Data.HighBeamHeadLight()),
			},
		}
		data = append(data, endpointData)
	}

	// response to browser
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   data,
	})
}

func Descrypt_mam_response(mam_message string, private_key string) *OBD.OBD2Meta {

	endpointSerial := endpoint_deserializer.Endpoint_Msg_flatbuffer_deserialize(mam_message)
	ciphertext := endpointSerial.DataBytes()
	// hmac := endpointSerial.HmacBytes()
	iv := endpointSerial.IvBytes()

	compressData := endpoint_CBCDecrypter.Endpoint_CBCDecrypter(string(ciphertext), private_key, string(iv))

	obd2MetaFlatbuffer, _ := cbrotli.Decode([]byte(compressData))

	obd2Meta := OBD.GetRootAsOBD2Meta(obd2MetaFlatbuffer, 0)
	log.Println("Descrypt mam response successfully")
	return obd2Meta
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
