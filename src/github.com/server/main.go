package main

import (
	"flag"
	"net/http"
	"github.com/gorilla/websocket"
	"net/url"
	"log"
	//"fmt"
)

type msgStruct struct {
	UserId int
	Content []byte
}

type userConnect struct {
	UserId int
	Output chan []byte //output channel for a user
	Input chan []byte //input channel from user
	Sid []byte      //session id
}

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
var upgrader = websocket.Upgrader{}
//массив структур с id пользователя и исходящим каналом к нему
var userConnects = make([]userConnect, 0, 10)
var count = 0

/**
 * регистрируем пользователя, выдавая коннект
 */
func register(userId int, input chan []byte) userConnect {
	var userConnect userConnect
	log.Printf("register user: %d", userId)

	userConnect.Output = make(chan []byte, 1024)
	userConnect.UserId = userId
	userConnect.Input = input

	userConnects = append(userConnects, userConnect)

	return userConnect
}

/**
 * роутим сообщение нужным пользователям
 */
func route(msg msgStruct) {
	for _, userConnect := range userConnects {
		if msg.UserId != userConnect.UserId {
			userConnect.Output <- msg.Content
		}
	}
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

	//превращаем Reader в chan
	userInput := make(chan []byte, 1024)
	go func(ch chan []byte) {
		for {
			if _, message, err := conn.ReadMessage(); err == nil {
				log.Printf("recv: %s", message)
				ch <- message
			} else {
				log.Print(err)
				conn.Close()
				break
			}
		}
	} (userInput)

	userConnect := register(userId, userInput)

	for {
		select {
		//если прочитали из сокета от юзера - роутим сообщение
		case bytes := <- userConnect.Input:
			msg := msgStruct{UserId: userId, Content: bytes}
			log.Printf("user %d print: %s", msg.UserId, msg.Content)
		//асинхронно рассылаем сообщение пользователям
			go route(msg)
		//если прочитали что-то из входящего канала юзеру - отправляем в сокет
		case bytes := <- userConnect.Output:
			log.Printf("send message \"%s\" for user %d", bytes, userConnect.UserId)
			if err = conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				log.Print(err)
				conn.Close()
				break
			}
		}
	}
}

func main() {
	upgrader.CheckOrigin = checkSameOrigin
	http.HandleFunc("/", handle)
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