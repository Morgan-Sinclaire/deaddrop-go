package read

import (
	"fmt"
	"log"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"github.com/Morgan-Sinclaire/deaddrop-go/session"
	"github.com/Morgan-Sinclaire/deaddrop-go/logging"
)

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

	messages := db.GetMessagesForUser(user)
	for _, message := range messages {
		fmt.Println(session.Decrypt(message))
		logging.LogMessage("Read messages for user " + user + "")
	}
}
