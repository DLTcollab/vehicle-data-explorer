package test

import (
	"testing"

	"biilabs.com/mam-data-explorer/models/endpoint_deserializer"
	"github.com/stretchr/testify/assert"
)

func Test_endpoint_deserializer(t *testing.T) {
	var test_timestamp uint64 = 1589274915
	var test_hmac = []byte{37, 39, 48, 30, 53, 63, 45, 61, 87, 60, 60, 90, 63, 63, 13, 94, 53, 93, 52, 7, 19, 38, 15, 31, 29, 32, 97, 2, 47, 67, 59, 62}
	var test_payload = []byte{
		99, 44, 121, 217, 149, 161, 127, 33, 133, 77, 125, 156, 53, 53, 248, 95, 57, 196, 141, 90, 121, 158,
		133, 218, 153, 153, 24, 84, 32, 245, 68, 131, 33, 189, 93, 182, 94, 220, 215, 227, 42, 85, 127, 95,
		138, 119, 190, 196, 60, 75, 30, 181, 233, 164, 143, 130, 61, 167, 214, 93, 156, 26, 225, 189, 216, 62,
		116, 54, 26, 75, 26, 68, 160, 153, 163, 43, 17, 97, 239, 77, 172, 13, 0, 149, 177, 145, 24, 239,
		57, 238, 76, 213, 9, 45, 147, 225, 107, 7, 23, 134, 82, 49, 202, 243, 203, 110, 30, 220, 207, 13,
		41, 124, 26, 43, 17, 204, 188, 41, 187, 245, 24, 7, 203, 33, 53, 94, 2, 160, 101, 25, 38, 183,
		75, 241, 170, 22, 95, 200, 242, 46, 213, 27, 170, 240, 70, 188, 188, 2, 229, 119, 248, 253, 126, 195,
		30, 179, 33, 32, 84, 134, 58, 122, 61, 133, 107, 232, 155, 202, 176, 141, 249, 134, 168, 163, 118, 238,
		95, 50, 240, 69, 169, 232, 66, 39, 171, 97, 219, 204, 129, 47, 82, 187, 169, 144, 64, 21, 120, 219,
		223, 40, 104, 216, 174, 16, 124, 36, 254, 219, 86, 239, 32, 255, 215, 99, 39, 131, 196, 2, 79, 69,
		49, 162, 1, 218, 50, 65, 239, 170, 29, 207, 210, 133, 167, 129, 150, 35, 165, 148, 255, 252, 131, 31,
		251, 91, 130, 34, 222, 70, 36, 45, 140, 85, 207, 141, 48, 1, 206, 31, 171, 235, 238, 126, 113}

	var test_iv = []byte{164, 3, 98, 193, 52, 162, 107, 252, 184, 42, 74, 225, 157, 26, 88, 72}

	var test_bytes = []byte{164, 3, 98, 193, 52, 162, 107, 252, 184, 42, 74, 225, 157, 26, 88, 72, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 49, 53, 56, 57, 50, 55, 52, 57, 49, 53, 37, 39, 48, 30, 53, 63, 45, 61, 87, 60, 60, 90, 63, 63, 13, 94, 53, 93, 52, 7, 19, 38, 15, 31, 29, 32, 97, 2, 47, 67, 59, 62, 48, 48, 48, 48, 48, 48, 48, 50, 54, 51, 99, 44, 121, 217, 149, 161, 127, 33, 133, 77, 125, 156, 53, 53, 248, 95, 57, 196, 141, 90, 121, 158, 133, 218, 153, 153, 24, 84, 32, 245, 68, 131, 33, 189, 93, 182, 94, 220, 215, 227, 42, 85, 127, 95, 138, 119, 190, 196, 60, 75, 30, 181, 233, 164, 143, 130, 61, 167, 214, 93, 156, 26, 225, 189, 216, 62, 116, 54, 26, 75, 26, 68, 160, 153, 163, 43, 17, 97, 239, 77, 172, 13, 0, 149, 177, 145, 24, 239, 57, 238, 76, 213, 9, 45, 147, 225, 107, 7, 23, 134, 82, 49, 202, 243, 203, 110, 30, 220, 207, 13, 41, 124, 26, 43, 17, 204, 188, 41, 187, 245, 24, 7, 203, 33, 53, 94, 2, 160, 101, 25, 38, 183, 75, 241, 170, 22, 95, 200, 242, 46, 213, 27, 170, 240, 70, 188, 188, 2, 229, 119, 248, 253, 126, 195, 30, 179, 33, 32, 84, 134, 58, 122, 61, 133, 107, 232, 155, 202, 176, 141, 249, 134, 168, 163, 118, 238, 95, 50, 240, 69, 169, 232, 66, 39, 171, 97, 219, 204, 129, 47, 82, 187, 169, 144, 64, 21, 120, 219, 223, 40, 104, 216, 174, 16, 124, 36, 254, 219, 86, 239, 32, 255, 215, 99, 39, 131, 196, 2, 79, 69, 49, 162, 1, 218, 50, 65, 239, 170, 29, 207, 210, 133, 167, 129, 150, 35, 165, 148, 255, 252, 131, 31, 251, 91, 130, 34, 222, 70, 36, 45, 140, 85, 207, 141, 48, 1, 206, 31, 171, 235, 238, 126, 113}
	var payload_len int = 263

	var s = string(test_bytes)
	var endpoint_serial = endpoint_deserializer.Endpoint_deserializer(s)

	assert.Equal(t, []byte(endpoint_serial.IV), test_iv)
	assert.Equal(t, endpoint_serial.Timestamp, test_timestamp)
	assert.Equal(t, []byte(endpoint_serial.Hmac), test_hmac)
	assert.Equal(t, endpoint_serial.Ciphertext_len, payload_len)
	assert.Equal(t, []byte(endpoint_serial.Ciphertext), test_payload)
}