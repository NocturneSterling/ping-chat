package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/scrypt"
)

func deriveKey(pass string) []byte {
	hash := sha256.Sum256([]byte(pass))
	key, _ := scrypt.Key([]byte(pass), hash[:8], 32768, 8, 1, 32)
	return key
}

func encrypt(text, pass string) []byte {
	hash := sha256.Sum256([]byte(pass))
	key := deriveKey(pass)
	aead, _ := chacha20poly1305.New(key)
	nonce := make([]byte, aead.NonceSize())
	rand.Read(nonce)
	ct := aead.Seal(nonce, nonce, []byte(text), nil)
	return append(hash[:], ct...)
}

func extractHash(data []byte) []byte {
	return data[:32]
}

func decrypt(data []byte, pass string) string {
	key := deriveKey(pass)
	aead, _ := chacha20poly1305.New(key)
	ct := data[32:]
	ns := aead.NonceSize()
	plain, err := aead.Open(nil, ct[:ns], ct[ns:], nil)
	if err != nil {
		return ""
	}
	return string(plain)
}

func main() {
	encrypted := encrypt("verysecuretext", "password123")
	fmt.Println(encrypted)
	fmt.Println(decrypt(encrypted, "password123"))
	fmt.Println(extractHash(encrypted))
}
