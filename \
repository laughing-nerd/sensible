package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

const (
	saltSize  = 16
	nonceSize = 12
	keySize   = 32 // 256 bits for AES-256
)

// deriveKey generates a 256-bit key using Argon2id
func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, keySize)
}

// Encrypt encrypts plaintext using a password and returns final blob [salt | nonce | ciphertext]
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	final := append(salt, nonce...)
	final = append(final, ciphertext...)
	return final, nil
}

// Decrypt decrypts the blob [salt | nonce | ciphertext] using the password
func Decrypt(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < saltSize+nonceSize {
		return nil, errors.New("invalid encrypted data")
	}

	salt := encrypted[:saltSize]
	nonce := encrypted[saltSize : saltSize+nonceSize]
	ciphertext := encrypted[saltSize+nonceSize:]

	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AskPassword interactively reads password from user
func AskPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	pwBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(pwBytes), err
}
