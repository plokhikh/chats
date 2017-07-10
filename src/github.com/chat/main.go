package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/websocket"
	"net/url"
	"log"
	"encoding/json"
	"strconv"
	"fmt"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{}
//массив структур с id пользователя и исходящим каналом к нему
var online = make([]UserOnline, 0, 10)
var count = 0
var input = make(chan RawPkg)

/** регистрируем пользователя, выдавая коннект */
func enter(userId int, conn *websocket.Conn) (UserOnline) {
	var userOnline UserOnline
	log.Printf("user entered: %d", userId)

	userOnline.Output = conn
	userOnline.UserId = userId

	raiseEvent(
		eventNames[EventNameUserEntered],
		Source{Type: SourceTypeUser, Guid: strconv.Itoa(userId), Name: fmt.Sprintf("user %d", userId)},
		nil,
	)

	online = append(online, userOnline)

	return userOnline
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

	userOnline := enter(userId, conn)

	for {
		if _, message, err := conn.ReadMessage(); err == nil {
			log.Printf("recv: %s", message)
			input <- RawPkg{User: userOnline, Message: message}
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
		//@todo check that's real command
		command := Command{}

		if err := json.Unmarshal([]byte(rawPkg.Message), &command); err!=nil {
			log.Panic(err)
		} else {
			command.Source.Type = SourceTypeUser
			command.Source.Guid = strconv.Itoa(rawPkg.User.UserId)

			//если у команды есть какой-то результат работы - отправляем его обратно в сокет
			if result := commandHandler(command); result != nil {
				encoded, _ := json.Marshal(result)
				log.Printf("Send message \"%s\" to user with guid \"%d\"", encoded, rawPkg.User.UserId)

				rawPkg.User.Output.WriteMessage(websocket.TextMessage, encoded)
			}
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