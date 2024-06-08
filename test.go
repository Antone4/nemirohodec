package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		Printfln("%v", r.Header.Get("Origin"))
		return strings.Contains(r.Header.Get("Origin"), ("http://" + thisPage))
	},
}
var (
	mu      sync.Mutex
	num     int
	clients = make(map[*websocket.Conn]bool)
)

type Comment struct {
	Username string    `json:"username"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
}

func HandleTest() {
	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/test.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templateData, err := os.ReadFile("templates/defines.html")
		if err != nil {
			panic(err)
		}
		tmpl = template.Must(tmpl.Parse(string(templateData)))
		session, _ := store.Get(r, "authdata")
		interfaceEmail, authorized := session.Values["email"]
		email := ""
		if authorized {
			email = interfaceEmail.(string)
		}

		err = tmpl.Execute(w, &Context{r, "hi", "", &tovars, email, authorized})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}).Methods("GET")

	router.HandleFunc("/test", AuthHandler).Methods("POST")
	router.HandleFunc("/deleteCookie", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "authdata")
		session.Options.MaxAge = -1
		session.Save(r, w)
		w.Write([]byte("Cookie has been deleted"))
	})

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		connWebSocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			Printfln("Error1: %v", err.Error())
			return
		}
		mu.Lock()
		clients[connWebSocket] = true
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(clients, connWebSocket)
			mu.Unlock()
			connWebSocket.Close()
		}()
		// for {
		// 	var comment Comment
		// 	err := connWebSocket.ReadJSON(&comment)
		// 	if err != nil {
		// 		Printfln("Error2: %v", err.Error())
		// 		break
		// 	}
		// 	comment.Username = "HUI"
		// 	err = connWebSocket.WriteJSON(comment)
		// 	if err != nil {
		// 		Printfln("Error3: %v", err.Error())
		// 		break
		// 	}
		// }
		for {
			_, message, err := connWebSocket.ReadMessage()
			if err != nil {
				Printfln("Read error: %v", err)
				break
			}
			Printfln("%v", string(message))
			mu.Lock()
			if string(message) == "increment" {
				num++
			} else if string(message) == "decrement" {
				num--
			}
			currentNum := num
			mu.Unlock()
			Printfln("%v", num)
			broadcastMessage(strconv.Itoa(currentNum))
		}
	}).Methods("GET")
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	Printfln("SUBMIT")
}

func broadcastMessage(message string) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			Printfln("Write error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
