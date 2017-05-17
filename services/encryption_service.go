package services

import (
	"crypto/md5"
	"encoding/hex"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	b32 "encoding/base32"
	"math/rand"
)

var encryptionKey = os.Getenv("ENCRYPTION_KEY")
var commonIV = []byte(os.Getenv("COMMON_IV"))

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

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

	sEnc := b32.StdEncoding.EncodeToString(ciphertext)
	fmt.Printf("\nFinal token: %s", sEnc)
	return sEnc
}

func DecodeToken(input string) string {

	fmt.Printf("\nInput: %s", input)

	sDec, _ := b32.StdEncoding.DecodeString(input)
	c, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		panic(err)
	}
	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	decryptedToken := make([]byte, len(sDec))
	cfbdec.XORKeyStream(decryptedToken, []byte(sDec))
	fmt.Printf("\nDecrypting: %x=>%s", sDec, decryptedToken)
	return string(decryptedToken)

}

func RandomPassword() string {
	return RandStringBytes(10)
}

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}





