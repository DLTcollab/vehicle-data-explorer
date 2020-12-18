package test

import (
	"testing"

	"github.com/DLTcollab/vehicle-data-explorer/OBD"
	controller "github.com/DLTcollab/vehicle-data-explorer/controllers"
	"github.com/stretchr/testify/assert"
)

func Test_Descrypt_mam_response(t *testing.T) {
	testMAMMessage := "EAAAAAAACgAQAAwACAAEAAoAAAAMAAAAvAAAAMwAAACwAAAAXu3axt+gxoFbKbYrSIRuP0Sv7IOgKgycT0WcUIJlUrMDOLE/GGDfhOdnmk8MVUrFEKxnGxPv9ZEq3u7ndbBNyC/xkTT2obqLRJRa957Kio1HVCqo16gibne7y+nDDeAvPW4BehAPam/6/YauoNMjp6xT+77sCnpLN6kd829tePKwNqk3mYtQPvULYgEQ+3tFh1LH9b7Zw6BR7m73wG2bEdQluWzD/31wpX/nyGf/RSUQAAAAAAAAAAAAAAAAAAAAAAAAACAAAAC5olhadUJd4z4Vh85AN5ktj/VCL+fH/Itd0xA0ksEl9w=="
	testPrivateKey := "LLHRCBHHYWKAGXMYCEKJIPBATQZPBQIE"
	testDeviceID := "PWFJOIOZGKUIOBVY"

	obd2Meta := controller.Descrypt_mam_response(testMAMMessage, testPrivateKey)
	obd2Data := new(OBD.OBD2_data)
	obd2Meta.Data(obd2Data)
	assert.Equal(t, string(obd2Meta.DeviceID()), testDeviceID)
	assert.Equal(t, string(obd2Data.Vin()), testDeviceID)
}
