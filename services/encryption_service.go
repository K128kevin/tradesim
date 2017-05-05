package services

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	b64 "encoding/base64"
)

var encryptionKey = "g4WLBBlAgMJgZkmH13szM3OenpMuPw2S"
var commonIV = []byte("FaOmz3fR3Rjeo1UeAUNDhd5ou1lmhDCD")

func HashString(input string) string {
	hasher := md5.New()
    hasher.Write([]byte(input))
    return hex.EncodeToString(hasher.Sum(nil))
}

func CreateToken(input string) string {
	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(input))
	cfb.XORKeyStream(ciphertext, []byte(input))
	fmt.Printf("\nEncrypting: %s=>%x", []byte(input), ciphertext)

	sEnc := b64.StdEncoding.EncodeToString(ciphertext)
	return sEnc
}

func DecodeToken(input string) string {

	sDec, _ := b64.StdEncoding.DecodeString(input)
	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		panic(err)
	}
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	decryptedToken := make([]byte, len(sDec))
	cfbdec.XORKeyStream(decryptedToken, []byte(sDec))
	fmt.Printf("Decrypting: %x=>%s", sDec, decryptedToken)
	return string(decryptedToken)

}





