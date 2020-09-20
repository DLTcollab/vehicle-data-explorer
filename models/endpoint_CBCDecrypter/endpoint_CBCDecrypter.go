package endpoint_CBCDecrypter

import (
	"crypto/aes"
	"crypto/cipher"
)

func Unpadding(src []byte) []byte {
	n := len(src)
	if n == 0 {
		return src
	}
	paddingNum := int(src[n-1])
	return src[:n-paddingNum]
}

func Endpoint_CBCDecrypter(ciphertext_str string, key string, iv_str string, timestamp uint64) string {

	var ciphertext = []byte(ciphertext_str)

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	var iv = []byte(iv_str)

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = Unpadding(ciphertext)
	return string(ciphertext)
}
