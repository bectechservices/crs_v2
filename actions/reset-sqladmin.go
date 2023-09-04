package actions

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var PasswordEncryptionKey = "9z$C&F)J@NcRfUjXn2r5u7x!A%D*G-Ka"

func EncryptPassword(text []byte) []byte {
	c, err := aes.NewCipher([]byte(PasswordEncryptionKey))
	if err != nil {
		fmt.Println(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	return gcm.Seal(nonce, nonce, text, nil)
}

func DecryptPassword(text []byte) string {
	c, err := aes.NewCipher([]byte(PasswordEncryptionKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		panic(err)
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		panic(err)
	}

	nonce, text := text[:nonceSize], text[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, text, nil)
	if err != nil {
		panic(err)
	}
	return string(plaintext)
}

func ReadPasswordFromFile() []byte {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	password, err := ioutil.ReadFile(filepath.Join(exPath, "becencrpt.txt"))

	if err != nil {
		panic(err)
		//	fmt.Println(password)
	}
	return password
}
