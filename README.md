# deaddrop-go

A deaddrop utility written in Go. Put files in a database behind a password to be retrieved at a later date.

This is a part of the University of Wyoming's Secure Software Design Course (Spring 2023). This is the base repository to be forked and updated for various assignments. Alternative language versions are available in:
- [Javascript](https://github.com/andey-robins/deaddrop-js)
- [Rust](https://github.com/andey-robins/deaddrop-rs)

## Versioning

`deaddrop-go` is built with:
- go version go1.19.4 linux/amd64

## Usage

`go run main.go --help` for instructions

Then run `go run main.go -new -user <username here>` and you will be prompted to create the initial password.

## Database

Data gets stored into the local database file dd.db. This file will not by synched to git repos. Delete this file if you don't set up a user properly on the first go

## Logging Strategy

2/12/23: We created a "logging" folder containing "logging.go", which in turn has a simple logging function which writes to "logs.txt". This function is used in new.go, read.go, and send.go to indicate what time relevant actions have happened. We kept track of six possible actions:

1) Creating a new user
2) Reading a user's messages
3) Trying to read a user's messages but giving the wrong password
4) Trying to read a user's messages but not giving a valid username
5) Send a message to a user
6) Trying to send a message to a user but not giving a valid username

In each case, we recorded the time of the event, the username in question, and which of the 6 events it was (in words). We think this is not overly verbose as it is one line and records 3 essential pieces of information, but not anything else. For instance, we wouldn't want to record the content of the messages, as this would defeat the point of confidentiality and reintroduce the problem we dealt with in the "Mitigation" section below.

## Mitigation

2/12/23: We pointed out previously that "dd.db" could be opened and was unencrypted, hence could be read/rewritten arbitrarily. I worked with Long on this: I wanted to password-protect dd.db, which is simple to do on the Mac OS UI using Disk Utility, but we had no clue how to do in Go. He wanted to encrypt dd.db instead. My own critique of this was that an encrypted database could still have data deleted, but some sort of system of admin permissions could prevent this. In any case, we tried to figure these out in parallel for a few hours. I'd been mainly trying to follow along with this source, but I was not able to install and run the library successfully (I suspect it's deprecated, but I can't prove this): https://ostechnix.com/cryptogo-easy-way-encrypt-password-protect-files/ But he succeeded first by modifying the code from: https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes. The code for the Encrypt/Decrypt functions in session.go is literally just copied from what he got.

## Message Authentication

3/26/23: First, we made deaddrop into a live drop, i.e. we changed the "send" mode so that a "-from" field needs to be given, and the sender needs to authenticate themselves. That is we changed the syntax:

go run main.go –send –to <username>

into:

go run main.go –send –to <recipient> -from <sender>

Then we added a MAC to all messages as they are added, and provided a means to compare the hashes. If they fail to compare, the user is warned whenever they access the message, otherwise they are told that it checks out.
