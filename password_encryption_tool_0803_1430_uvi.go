// 代码生成时间: 2025-08-03 14:30:13
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "io"
    "log"
    "os"
)

// Utility functions

// EncryptPassword encrypts the given password using AES encryption
func EncryptPassword(password string) (string, error) {
    key := []byte("your-encryption-key") // Ensure this key is kept secure
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    // Pad the password to be a multiple of the block size
    blockSize := block.BlockSize()
    paddedPassword := PKCS5Padding([]byte(password), blockSize)

    // Generate a random initialization vector
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    // Encrypt the padded password
    mode := cipher.NewCBCEncrypter(block, iv)
    encryptedPassword := make([]byte, len(paddedPassword))
    mode.CryptBlocks(encryptedPassword, paddedPassword)

    // Encode the initialization vector and encrypted password in base64
    encodedIV := base64.StdEncoding.EncodeToString(iv)
    encodedPassword := base64.StdEncoding.EncodeToString(encryptedPassword)
    return encodedIV + ":" + encodedPassword, nil
}

// DecryptPassword decrypts the encrypted password
func DecryptPassword(encryptedPassword string) (string, error) {
    key := []byte("your-encryption-key") // Ensure this key is kept secure
    parts := split(encryptedPassword, ':')
    if len(parts) != 2 {
        return "", fmt.Errorf("invalid format")
    }

    encodedIV := parts[0]
    encodedPassword := parts[1]
    iv, err := base64.StdEncoding.DecodeString(encodedIV)
    if err != nil {
        return "", err
    }
    encryptedPass, err := base64.StdEncoding.DecodeString(encodedPassword)
    if err != nil {
        return "", err
    }

    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    // Decrypt the password
    mode := cipher.NewCBCDecrypter(block, iv)
    decryptedPassword := make([]byte, len(encryptedPass))
    mode.CryptBlocks(decryptedPassword, encryptedPass)

    // Unpad the decrypted password
    unpaddedPassword := PKCS5UnPadding(decryptedPassword)
    return string(unpaddedPassword), nil
}

// PKCS5Padding pads the data to be a multiple of the block size
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
    padding := blockSize - len(ciphertext)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(ciphertext, padtext...)
}

// PKCS5UnPadding removes the PKCS5 padding from the data
func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

// split splits a string by a delimiter and returns a slice of the parts
func split(s, delimiter string) []string {
    return strings.Split(s, delimiter)
}

func main() {
    // Example usage
    password := "mysecretpassword"
    encryptedPass, err := EncryptPassword(password)
    if err != nil {
        log.Fatalf("Error encrypting password: %v", err)
    }
    fmt.Printf("Encrypted Password: %s
", encryptedPass)

    decryptedPass, err := DecryptPassword(encryptedPass)
    if err != nil {
        log.Fatalf("Error decrypting password: %v", err)
    }
    fmt.Printf("Decrypted Password: %s
", decryptedPass)
}
