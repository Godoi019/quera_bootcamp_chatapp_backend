package main

import (
	"log"
	"os"

	"github.com/Hossara/quera_bootcamp_chatapp_backend/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
