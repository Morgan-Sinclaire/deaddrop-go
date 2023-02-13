package db

import (
	"log"
)

// GetMessagesForUser assumes that a user has already been
// authenticated through a call to session.Authenticate(user)
// and then returns all the messages stored for that user
func GetMessagesForUser(user string) [][]byte {
	database := Connect().Db

	rows, err := database.Query(`
		SELECT (data) FROM Messages
		WHERE recipient = (
			SELECT id FROM Users WHERE user = ?
		)
		`, user)
	if err != nil {
		log.Fatalf("err: %v", err)
	}
	defer rows.Close()

	// marshall rows into an array
	messages := make([][]byte, 0)
	for rows.Next() {
		var message []byte
		err := rows.Scan(&message)
		if err != nil {
			log.Fatalf("unable to scan row")
		}
		messages = append(messages, message)
	}
	return messages
}

// saveMessage will process the transaction to place a message
// into the database
func SaveMessage(message, recipient string) {
	database := Connect().Db

	database.Exec(`
		INSERT INTO Messages (recipient, data)
		VALUES (
			(SELECT id FROM Users WHERE user = ?),
			?
		);
	`, recipient, message)
}
