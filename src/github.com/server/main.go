package main

import (
	_ "fmt"
	_ "net"
	_ "os"
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"net/url"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{}

func handle(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		if mt, message, err := conn.ReadMessage(); err != nil {
				panic(err)
		} else {
			log.Printf("recv: %s", message)
			if err = conn.WriteMessage(mt, message); err != nil {
				panic(err)
			}
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
}

func main() {
	upgrader.CheckOrigin = checkSameOrigin
	http.HandleFunc("/echo", handle)
	http.HandleFunc("/", home)
	panic(http.ListenAndServe(*addr, nil))
}

func checkSameOrigin(r *http.Request) bool {
	//warning! security vulnerable! temporary turn off origin checking
	return true
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	return u.Host == r.Host
}