package db

import (
	"log"
)

// GetMessagesForUser assumes that a user has already been
// authenticated through a call to session.Authenticate(user)
// and then returns all the messages stored for that user
func GetMessagesForUser(user string) ([][]byte,[]string,[]string) {
	// []([]byte, string, string)

	database := Connect().Db

	rows, err := database.Query(`
		SELECT data, user, Messages.hash
		FROM Messages INNER JOIN Users ON Users.id=Messages.sender
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
	senders := make([]string, 0)
	hashes := make([]string, 0)
	for rows.Next() {
		var message []byte
		var sender string
		var hash string
		err := rows.Scan(&message, &sender, &hash)
		if err != nil {
			log.Fatalf("unable to scan row")
		}
		messages = append(messages, message)
		senders = append(senders, sender)
		hashes = append(hashes, hash)

	}
	return messages,senders,hashes
}

// saveMessage will process the transaction to place a message
// into the database
func SaveMessage(message, recipient string, sender string, hash string) {
	database := Connect().Db

	database.Exec(`
		INSERT INTO Messages (recipient, sender, data, hash)
		VALUES (
			(SELECT id FROM Users WHERE user = ?),
			(SELECT id FROM Users WHERE user = ?),
			?,
			?
		);
	`, recipient, sender, message, string(hash))
}
