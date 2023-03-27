package read

import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"github.com/Morgan-Sinclaire/deaddrop-go/session"
	"github.com/Morgan-Sinclaire/deaddrop-go/logging"
)

// ReadMessages takes a username, prompts that user for their password,
// and (if authentication succeeds) lets that user read their messages
func ReadMessages(user string) {
	if !db.UserExists(user) {
		logging.LogMessage("Tried to read messages for nonexistent user " + user + "")
		log.Fatalf("User not recognized")
	}

	err := session.Authenticate(user)
	if err != nil {
		logging.LogMessage("Tried to read messages for user " + user + " with wrong password")
		log.Fatalf("Unable to authenticate user")
	}

	messages,senders,hashes := db.GetMessagesForUser(user)
	for i, message := range messages {
		sender := senders[i]
		hash := hashes[i]
		plaintext := session.Decrypt(message)
		fmt.Print("message: " + plaintext)
		fmt.Println("sender: " + sender)
		logging.LogMessage("Read messages for user " + user + "")
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plaintext + session.KEY + sender))
		fmt.Println([]byte(plaintext + session.KEY + sender))
		// fmt.Println("hash: " + hash)
		// bc, _ := bcrypt.GenerateFromPassword([]byte(plaintext + session.KEY + sender), bcrypt.DefaultCost)
		// fmt.Println("bcrypt: " + string(bc))
		if err != nil {
			fmt.Println("The integrity of this message could not be verified.")
		} else {
			fmt.Println("The integrity of this message has been verified.")
		}
		fmt.Println("")
	}
}
