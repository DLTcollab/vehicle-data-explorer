package test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	controller "github.com/DLTcollab/vehicle-data-explorer/controllers"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_CBCDecrypter"
	"github.com/DLTcollab/vehicle-data-explorer/models/endpoint_deserializer"
	"github.com/stretchr/testify/assert"
)

func Test_descrypt_endpoint_message(t *testing.T) {
	mam_message := "123456789012345600000000001600653076AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA0000001376fb8a5d80f6e2ee8450fc0e32f6d950ce2e4b6432000d7b815f0244be37bbcdfc53e3523bbbc9ed878625ca38447d0c69e64be472ace636e9cc72323493182bf2cb10acee566635d706364ef4c6f0c0ed9edf389be4552c18f5a04ae1552037ecf247a667c94d2cfdd38e1d308a9b4a077a52e0151b218772a8d29da1478d5ce17e13b507207a1b342db8020afcad5ccf5e080236ec1914ad5d78fdb11256dee52b424f50ab1324c03788a25b10eb2f7bba1a0e414ec83b172a35ac3e9a4423a5879cdfa15bc8b55a46048cef9c7f6dea611df7d5bf96748cd3594147d89addbe12d0b56323a4fd7b70ae72f582a77fd2552e80e69147ad56bd369ab85d81ef11dded28b26a818be2e14ffa5a9ec112efa7978cec2ee683d0db2d449282b02106efce2a8db02ece41fc0c0a17a203169df310cd8db79c56e2fbb55115454e66eb34f32896d335623041c2092b397cf41e8158bb840612066cace26d5950a1e69bc16876c5c6c029bf7842a7c76d35b6b450e26f0e787e27fb5ecb6b4721d9dc52158422e73c0f48db979a179e247f42e02081fcaa6aaef12dd196bf4021d947773fc33f0d9c614e8ca01924d56385b7b5d0d2ad53ef74ac83bedb38083b08d453f360498f31b9792ed54c9d6b137dcaed534d5ea92567ad4f6719bec39640ef28a08cee450b2efe2539c0acc48d7dced9848dab3d8e2090482f6ca431aadeb5d02ad3e659c3dea8de04b9fc02acf6f17d15361907f464073f11ebeac727321d2858c2d85f68b1c17d4684f44720b2805df137beec904a13ffa310ce8c87d1f97f282759dbc70d306886e3a5cc46046b3f52a657b8a4d9f04ffef1e1bf59dfdbea9d13217aa9086139412840a255d329c6ef1e347f48328078b006d29c0cc4ea79a7ca698eb3ec442979d8af61b233ddf808f0bcd81bcf23c6fe0d689ae6dc660ec4472d677b61566232ef492efa9d0ea9"
	private_key := "LLHRCBHHYWKAGXMYCEKJIPBATQZPBQIE"
	origin_message := []byte(`{"timestamp":1600653148,"Device_id":"PWFJOIOZGKUIOBVY","data":{"vin":"BCKHQEUPEEQNMKGUJ","engine_load":25,"engine_coolant_temperature":51,"fuel_pressure":15,"engine_speed":57.72918,"vehicle_speed":10,"intake_air_temperature":10,"mass_air_flow":85,"fuel_tank_level_input":90,"absolute_barometric_pressure":32,"control_module_voltage":98.409615,"throttle_position":91,"ambient_air_temperature":82,"relative_accelerator_pedal_position":84,"engine_oil_temperature":97,"engine_fuel_rate":67.559395,"service_distance":71,"anti_lock_barking_active":94,"steering_wheel_angle":26,"position_of_doors":2,"right_left_turn_signal_light":81,"alternate_beam_head_light":79,"high_beam_head_light":66}}`)

	endpoint_serial := endpoint_deserializer.Endpoint_deserializer(mam_message)
	ciphertext, _ := hex.DecodeString(endpoint_serial.Ciphertext)
	plaintext := endpoint_CBCDecrypter.Endpoint_CBCDecrypter(string(ciphertext), private_key, endpoint_serial.IV)

	var endpoint_data, origin controller.Endpoint_obd2_data
	err := json.Unmarshal([]byte(plaintext), &endpoint_data)
	if err != nil {
		fmt.Println("error:", err)
	}

	json.Unmarshal(origin_message, &origin)
	assert.Equal(t, origin, endpoint_data)
}
