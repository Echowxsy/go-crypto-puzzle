package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil

}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func Encrypt(input, secret, salt []byte, pbkdf2Rounds int) ([]byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		_, err := rand.Read(salt)
		if err != nil {
			return nil, err
		}
	}

	if len(salt) != 32 {
		hash := sha256.Sum256(salt)
		salt = hash[:]
	}

	rounds := pbkdf2Rounds
	if rounds < 1 {
		rounds = 1
	}

	key := pbkdf2.Key(secret, salt, rounds, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, input, nil)
	version := byte(1)
	archive := make([]byte, 1+len(salt)+4+len(nonce)+len(ciphertext))
	archive[0] = version
	copy(archive[1:], salt)
	binary.BigEndian.PutUint32(archive[33:], uint32(rounds))
	copy(archive[37:], nonce)
	copy(archive[37+len(nonce):], ciphertext)

	return archive, nil
}

func Decrypt(input, secret []byte) ([]byte, error) {
	if len(input) < 37 {
		return nil, fmt.Errorf("invalid archive")
	}

	salt := input[1:33]
	rounds := binary.BigEndian.Uint32(input[33:37])

	nonce := input[37 : 37+12]
	ciphertext := input[37+12:]

	key := pbkdf2.Key(secret, salt, int(rounds), 32, sha256.New)

	fmt.Println("rounds: ", rounds)
	fmt.Println("salt: ", salt)
	fmt.Println("nonce: ", nonce)
	fmt.Println("ciphertext: ", ciphertext)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
