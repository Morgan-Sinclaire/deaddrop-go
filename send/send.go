package send

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"golang.org/x/crypto/bcrypt"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"github.com/Morgan-Sinclaire/deaddrop-go/logging"
	"github.com/Morgan-Sinclaire/deaddrop-go/session"
)

// SendMessage takes a destination username and a sending username.
// It will prompt for the sender's password, and if they authenticate
//successfully, it will prompt them for a message to send to the recipient.
func SendMessage(to string, from string) {
	if !db.UserExists(to) {
		logging.LogMessage("Tried to send message to nonexistent user " + to + "")
		log.Fatalf("Destination user does not exist")
	}

	if !db.UserExists(from) {
		logging.LogMessage("Tried to send messages from nonexistent user " + from + "")
		log.Fatalf("Sender username not recognized")
	}

	err := session.Authenticate(from)
	if err != nil {
		logging.LogMessage("Tried to send messages from user " + from + " with wrong password")
		log.Fatalf("Unable to authenticate user")
	}

	message := getUserMessage()

	hash, _ := bcrypt.GenerateFromPassword([]byte(message + session.KEY + from), bcrypt.DefaultCost)

	db.SaveMessage(session.Encrypt(message), to, from, string(hash))
	logging.LogMessage("Sent message to user " + to + " from " + from + "")
}

// getUserMessage prompts the user for the message to send
// and returns it
func getUserMessage() string {
	fmt.Println("Enter your message: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return text
}
