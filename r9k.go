package main

// https://github.com/bwmarrin/discordgo/blob/master/examples/dm_pingpong/main.go
// https://github.com/mattn/go-sqlite3/blob/master/_example/simple/simple.go

import (
	"fmt"
	"log"
	"github.com/bwmarrin/discordgo"
	"os/signal"
	"os"
	"syscall"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
var auth_token string
var db *sql.DB
var err error

func main() {
	// database columns
	// id INTEGER, message TEXT NOT NULL UNIQUE
	db, err = sql.Open("sqlite3", "data/messages.db")
	if err != nil {
		log.Fatal(err)
	}

	// AUTH TOKEN BELOW
	auth_token = ""
	// new discord session using bot token
	dg, err := discordgo.New("Bot " + auth_token)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// register messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageUpdate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	// cleanly close discord session
	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}

func messageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID != "" { // DECLARE CHANNEL ID
		return
	}
	// CHECK IF MESSAGE IS IN DATABASE
	rows, err := db.Query("select * from tblone where message=?", m.Content)
	defer rows.Close()
	if rows.Next() { // found match
		fmt.Println("Matching edited message found from",m.Author," (",m.Content,")")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
				log.Fatal(err)
		}

		_, err = s.ChannelMessageSend(channel.ID, constructMessage(m.Content))
		if err != nil {
				log.Fatal(err)
		}

	} else {         // no match
		_, err = db.Exec("insert into tblone(id, message) values(NULL, ?)", m.Content)
	if err != nil {
		log.Fatal(err)
	}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.ChannelID != "" { // DECLARE CHANNEL ID
		return
	}
	// CHECK IF MESSAGE IS IN DATABASE
	rows, err := db.Query("select * from tblone where message=?", m.Content)
	defer rows.Close()
	if rows.Next() { // found match
		fmt.Println("Matching message found from",m.Author," (",m.Content,")")
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		channel, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
				log.Fatal(err)
		}

		_, err = s.ChannelMessageSend(channel.ID, constructMessage(m.Content))
		if err != nil {
				log.Fatal(err)
		}

	} else {         // no match
		_, err = db.Exec("insert into tblone(id, message) values(NULL, ?)", m.Content)
	if err != nil {
		log.Fatal(err)
	}
	}
}

func constructMessage(m string) string {
	return ("Your message has failed to send as it matches a previous message.\n`"+m+"`")
}
