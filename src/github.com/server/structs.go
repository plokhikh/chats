package main

import "net/http"

type UserOnline struct {
	UserId int
	Output http.ResponseWriter //output writer for a user
	Sid []byte      //session id
}

/**
 * голая структура куда читается то, что присшло от юзера
 */
type RawPkg struct {
	UserId  int
	Message []byte
}

/**
 * базовая структура для всех сообщени от юзера
 */
type Pkg struct {
	Type string
	Ts int
	Data interface{}
}

/**
 * базовая структура для всех команд от пользователя
 */
type Command struct {
	*Pkg
	Action string
}

type CommandSendMessageData struct {
	To string
	Message string
}