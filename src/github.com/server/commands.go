package main

import (
	"log"
	"github.com/mitchellh/mapstructure"
)

func commandHandler(command Command) {
	switch name := command.Name; name {
	//@todo create enum for commands
	case "send-message":
		/** отдаем сырую мапу Data в функцию, которая сама разберет в нужную структуру*/
		commandSendMessage(command)
	default:
		log.Printf("Unknown command name %s", command.Name)
		log.Panic(command)
	}
}

func commandSendMessage(command Command) {
	var data MessageData
	if err := mapstructure.Decode(command.Data, &data); err != nil {
		log.Panic(err)
	} else {
		//@todo create event name enum
		raiseEvent("message-sended", command.Source, data)
	}
}