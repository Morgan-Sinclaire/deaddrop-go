package send

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"github.com/Morgan-Sinclaire/deaddrop-go/logging"
	"github.com/Morgan-Sinclaire/deaddrop-go/session"
)

// SendMessage takes a destination username and will
// prompt the user for a message to send to that user
func SendMessage(to string) {
	if !db.UserExists(to) {
		logging.LogMessage("Tried to send message to nonexistent user " + to + "")
		log.Fatalf("Destination user does not exist")
	}

	message := getUserMessage()

	db.SaveMessage(session.Encrypt(message), to)
	logging.LogMessage("Sent message to user " + to + "")
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
