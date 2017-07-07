package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/websocket"
	"net/url"
	"log"
	"encoding/json"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{}
//массив структур с id пользователя и исходящим каналом к нему
var online = make([]UserOnline, 0, 10)
var count = 0
var input = make(chan RawPkg)

/** регистрируем пользователя, выдавая коннект */
func register(userId int, response http.ResponseWriter) {
	var userConnect UserOnline
	log.Printf("register user: %d", userId)

	userConnect.Output = response
	userConnect.UserId = userId

	online = append(online, userConnect)
}

/**
 * обработка сокетного соединения браузером
 */
func handle(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Print(err)
	}

	//где-то здесь должна быть аутентификация
	count++
	userId := count

	register(userId, response)

	for {
		if _, message, err := conn.ReadMessage(); err == nil {
			log.Printf("recv: %s", message)
			input <- RawPkg{UserId: userId, Message: message}
			response.Write([]byte("sdfgsdgfdsg"))
		} else {
			log.Print(err)
			conn.Close()
			break
		}
	}
}

/** слушаем входящий канал и шлем куда надо */
func listenInput() {
	for {
		rawPkg := <- input
		command := Command{}

		if err := json.Unmarshal([]byte(rawPkg.Message), &command); err!=nil {
			log.Panic(err)
		} else {
			commandHandler(command)
		}
	}
}

func main() {
	upgrader.CheckOrigin = checkSameOrigin
	http.HandleFunc("/", handle)
	go listenInput()
	panic(http.ListenAndServe(*addr, nil))
}

func checkSameOrigin(r *http.Request) bool {
	origin := r.Header["Origin"]
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin[0])
	if err != nil {
		return false
	}
	//log.Println(u.Host, r.Host)
	//warning! security vulnerable! temporary turn off origin checking
	return true
	return u.Host == r.Host
}