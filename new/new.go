package new

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Morgan-Sinclaire/deaddrop-go/db"
	"github.com/Morgan-Sinclaire/deaddrop-go/session"
	"github.com/Morgan-Sinclaire/deaddrop-go/logging"
)

// Create a NewUser as authorized by the user 'user'
func NewUser(user string) {
	if !db.NoUsers() && !db.UserExists(user) {
		log.Fatalf("User not recognized")
	}

	err := session.Authenticate(user)
	if err != nil {
		log.Fatalf("Unable to authenticate user")
	}

	newUser := getNewUsername()
	newPassHash, err := session.GetPassword()
	if err != nil {
		log.Fatalf("Unable to get password hash")
	}

	err = db.SetUserPassHash(newUser, newPassHash)
	if err != nil {
		log.Fatalf("Unable to create new user")
	}

	logging.LogMessage("New user " + newUser + " created")
}

// getUserMessage prompts the user for the message to send
// and returns it
func getNewUsername() string {
	fmt.Println("Enter the username for the new user: ")
	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
	return strings.Trim(text, "\n\t ")
}
