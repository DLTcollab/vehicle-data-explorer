package endpoint_deserializer

import "strconv"

type Endpoint_serial struct {
	IV             string
	Timestamp      uint64
	Hmac           string
	Ciphertext_len int
	Ciphertext     string
}

func Endpoint_deserializer(serialize_msg string) Endpoint_serial {
	var endpoint_serial Endpoint_serial
	endpoint_serial.IV = serialize_msg[0:16]
	endpoint_serial.Timestamp, _ = strconv.ParseUint(serialize_msg[16:36], 10, 64)
	endpoint_serial.Hmac = serialize_msg[36:68]
	endpoint_serial.Ciphertext_len, _ = strconv.Atoi(serialize_msg[68:78])
	endpoint_serial.Ciphertext = serialize_msg[78:]
	return endpoint_serial
}
