package device

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Device struct {
	IMEI string `json:"imei"` // International Mobile Equipment Identity
	FSN  string `json:"fsn"`  // Serial number
	Name string `json:"name"` // Device name
}

func EncodeDevice(encodeDevice Device) string {
	JSONDevice, _ := json.Marshal(&encodeDevice)
	return string(JSONDevice)
}

func DecodeDevice(JSONDevice string) (*Device, error) {
	device := Device{}
	err := json.Unmarshal([]byte(JSONDevice), &device)

	if err != nil {
		return nil, err

	}
	return &device, nil
}

// Concatenate IMEI and FSN as identity
func GetDeviceHash(device Device) string {
	h := sha256.Sum256([]byte(device.IMEI + device.FSN))
	hexStr := hex.EncodeToString(h[:])
	return hexStr
}
