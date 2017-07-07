package main

import (
	"log"
	"github.com/mitchellh/mapstructure"
)

func commandHandler(command Command) {
	switch action := command.Action; action {
	case "send-message":
		/** отдаем сырую мапу Data в функцию, которая сама разберет в нужную структуру*/
		commandSendMessage(command.Data)
	default:
		log.Printf("Unknown command action %s", command.Action)
		log.Panic(command)
	}
}

func commandSendMessage(data interface{}) {
	var result CommandSendMessageData
	if err:=mapstructure.Decode(data, &result); err!=nil {
		log.Panic(err)
	} else {
		for _, userOnline := range online {
			log.Print(result.Message)
			userOnline.Output.Write([]byte(result.Message))
		}
	}
}