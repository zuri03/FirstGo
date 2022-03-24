package main

/*
*TODO:
	- Set up web socket
	- Set up server logging
	- Set up channels and signling for graceful shutdown
	- Try to simplify main
*/
import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	Server "github.com/zuri03/FirstGo/server"
)

func main() {

	cleanVariables := func(r rune) rune {
		if unicode.IsGraphic(r) {
			return r
		}
		return -1
	}

	data, err := os.ReadFile("./.env")
	if err != nil {
		log.Fatal(err)
		return
	}

	credentials := strings.Split(string(data), "\n")
	clientId := strings.Split(credentials[1], "=")
	clientSecret := strings.Split(credentials[0], "=")

	secret := strings.Map(cleanVariables, clientSecret[1])
	id := strings.Map(cleanVariables, clientId[1])
	os.Setenv(clientId[0], id)
	os.Setenv(clientSecret[0], secret)

	fmt.Println("====== INITIALIZING SERVER =====")
	file, err := os.OpenFile("./logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	Server.InitializeServer(file)

	for {
	}
}
