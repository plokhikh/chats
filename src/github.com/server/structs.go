package main

import (
	"github.com/gorilla/websocket"
)

type UserOnline struct {
	UserId int
	Output *websocket.Conn //output writer for a user
	Sid []byte      //session id
}

/**
 * голая структура куда читается то, что пришло от юзера
 */
type RawPkg struct {
	UserId  int
	Message []byte
}

/** источник события/команды */
type Source struct {
	Type SourceType `json:"type"`
	Guid string `json:"guid"`
	Name string `json:"name"`
}

/** базовая структура, некий "пакет", который может быть командой или событием */
type Pkg struct {
	Type PkgType `json:"type"`
	Ts int `json:"ts"` //unix timestamp with milliseconds
	Name string `json:"name"`
	Source Source `json:"source"`
	Data interface{} `json:"data"`
}

type Command struct {
	Pkg
}

type Event struct {
	Pkg
}

/**
 * специфические структуры данных для команд/событий
 */
type MessageData struct {
	To int `json:"to"`
	Message string `json:"message"`
}

type SendMessageCommand struct {
	Command
	Data MessageData `json:"data"`
}

/** send message events structure data */
type MessageSendedEvent struct {
	Event
	Data MessageData `json:"data"`
}
type UserEnteredEvent struct {
	Event
}