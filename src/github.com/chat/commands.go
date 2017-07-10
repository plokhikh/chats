package main

import (
	"log"
	"github.com/mitchellh/mapstructure"
	"time"
	"strconv"
	"fmt"
)

func commandHandler(command Command) (interface{}){
	switch name := command.Name; name {
	//@todo create enum for commands
	case commandNames[CommandNameSendMessage]:
		/** отдаем сырую мапу Data в функцию, которая сама разберет в нужную структуру*/
		return commandSendMessage(command)
	case commandNames[CommandNameGetUsersOnline]:
		return commandGetUsersOnline(command)
	default:
		log.Printf("Unknown command name %s", command.Name)
		log.Panic(command)
		return nil
	}
}

func commandSendMessage(command Command) (interface{}) {
	var data MessageData
	if err := mapstructure.Decode(command.Data, &data); err != nil {
		log.Panic(err)
	} else {
		raiseEvent(eventNames[EventNameMessageSended], command.Source, data)
	}

	return nil
}

func commandGetUsersOnline(command Command) (interface{}) {
	log.Print("got command get users online")
	var value UsersOnlineValue
	var sources []Source

	value.Type = PkgTypeValue
	value.Source = command.Source
	value.Name = valueNames[ValueNameUsersOnline]
	value.Ts = int(time.Now().UnixNano()/1000000)

	for _, userOnline := range online {
		sources = append(sources, Source{Type: SourceTypeUser, Guid: strconv.Itoa(userOnline.UserId), Name: fmt.Sprintf("user %d", userOnline.UserId)})
	}

	value.Data = sources

	return value
}