// 代码生成时间: 2025-09-08 00:01:37
// password_tool.go

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

// PasswordEncryptor is a struct that holds the key for encryption and decryption
type PasswordEncryptor struct {
	Key []byte
}

// NewPasswordEncryptor creates a new PasswordEncryptor with a given key
func NewPasswordEncryptor(key []byte) *PasswordEncryptor {
	return &PasswordEncryptor{Key: key}
}

// Encrypt encrypts a plaintext password using AES-256-GCM
func (p *PasswordEncryptor) Encrypt(plaintext string) (string, error) {
	if p.Key == nil || len(p.Key) != 32 {
		return "", errors.New("encryption key must be 32 bytes long")
	}

	block, err := aes.NewCipher(p.Key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// Decrypt decrypts an encrypted password using AES-256-GCM
func (p *PasswordEncryptor) Decrypt(ciphertext string) (string, error) {
	if p.Key == nil || len(p.Key) != 32 {
		return "", errors.New("encryption key must be 32 bytes long")
	}

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(p.Key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertextBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertextBytes[:nonceSize], ciphertextBytes[nonceSize:]
	return string(gcm.Open(nil, nonce, ciphertext, nil)), nil
}

func main() {
	// A 32-byte key for AES-256 encryption
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		log.Fatal(err)
	}

	encryptor := NewPasswordEncryptor(key)

	// Example usage
	plaintext := "mysecretpassword"
	encrypted, err := encryptor.Encrypt(plaintext)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Encrypted: %s
", encrypted)

	decrypted, err := encryptor.Decrypt(encrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Decrypted: %s
", decrypted)
}