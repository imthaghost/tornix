package main

import (
	"github.com/imthaghost/tornix"
	"io"
	"log"
)

func main() {
	manager := tornix.NewManager(10)

	session1, client1, err := manager.StartNewSession()
	if err != nil {
		log.Fatal(err)
	}

	manager.Sessions.Range(func(key, value interface{}) bool {
		log.Println(key, value)
		return true
	})

	resp, err := client1.Get("https://checkip.amazonaws.com")
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(session1)
	log.Println(string(body))

	// create a new session
	session2, client2, err := manager.StartNewSession()

	resp, err = client2.Get("https://checkip.amazonaws.com")
	if err != nil {
		log.Fatal(err)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(session2)
	log.Println(string(body))

	// end session
	log.Println("[+] Ending session")
	manager.EndSession(session1)
	manager.EndSession(session2)
}
