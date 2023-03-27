package session

import (
	"fmt"
	"log"
	"os"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
)

var KEY = "6b4cc601cd451c148d9aba4a3fdb544c"

func Encrypt(plainText string) (encryptedText string) {
	k, _ := hex.DecodeString(KEY)
	plain := []byte(plainText)

	block, err := aes.NewCipher(k)
	if err != nil {
		log.Fatal(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	encrypted := aesGCM.Seal(nonce, nonce, plain, nil)

	return string(encrypted)
}

func Decrypt(encryptedText []byte) (decryptedText string) {
	k, _ := hex.DecodeString(KEY)

	block, err := aes.NewCipher(k)
	if err != nil {
		log.Fatal(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err)
	}

	nonceSize := aesGCM.NonceSize()

	nonce, cipher := encryptedText[:nonceSize], encryptedText[nonceSize:]

	plain, err := aesGCM.Open(nil, nonce, cipher, nil)
	if err != nil {
		log.Fatal(err)
	}

	return string(plain)
}

// GetPassword will read in a password from stdin using the terminal
// no-echo utility ReadPassword. it will then salt and hash it with
// bcrypt
func GetPassword() (string, error) {
	pass, err := readPass()
	if err != nil {
		return "", err
	}

	return saltAndHash(pass)
}

// a nice wrapper to encapsulate the bcrypt generateFromPassword func
// to make salting and hashing easier
func saltAndHash(pass []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Authenticate takes in the username of a user and returns nil
// if the given password matches the user and an error otherwise
func Authenticate(user string) error {

	// bypass authentication if no users exist as this means
	// there is no data being stored and we can let the user create
	// a new user without auth
	if db.NoUsers() {
		return nil
	}

	pass, err := readPass()
	if err != nil {
		log.Fatalf("Error reading in password")
	}

	hash, err := db.GetUserPassHash(user)
	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(hash), pass)
}

// using the built in password read utility, get a password from stdin
func readPass() ([]byte, error) {
	fmt.Println("Password: ")
	pass, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return pass, nil
}
//
// func ProduceMAC(sender string, message string) (string, error) {
// 		hash, err := bcrypt.GenerateFromPassword([]byte(message + KEY + sender), bcrypt.DefaultCost)
// 		return string(hash), nil
// }
