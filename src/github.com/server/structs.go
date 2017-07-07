package main

import "net/http"

type pkg struct {
	UserId  int
	Message []byte
}

type userOnline struct {
	UserId int
	Output http.ResponseWriter //output writer for a user
	Sid []byte      //session id
}