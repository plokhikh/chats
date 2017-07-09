package main

import (
	"log"
	"time"
	"encoding/json"
	"github.com/gorilla/websocket"
)

func raiseEvent(name string, source Source, data interface{}) {
	var event Event

	event.Type = PkgTypeEvent
	event.Source = source
	event.Name = name
	event.Ts = int(time.Now().UnixNano()/1000000)

	switch name := event.Name; name {
	//@todo create enum for events
	case "message-sended":
	//assert that data interface exactly what eventMessageSended need
		if asserted, ok := data.(MessageData); ok {
			var messageSendedEvent MessageSendedEvent

			//could't copy Event field to MessageSendedEvent using MessageSendedEvent{Event}
			messageSendedEvent.Type = event.Type
			messageSendedEvent.Source = event.Source
			messageSendedEvent.Name = event.Name
			messageSendedEvent.Ts = event.Ts
			messageSendedEvent.Data = asserted

			eventMessageSended(messageSendedEvent)
		} else {
			log.Panic(data)
		}
	default:
		log.Printf("Unknown event name \"%s\"", event.Name)
		log.Panic(event)
	}
}

func eventMessageSended(event MessageSendedEvent) {
	for _, userOnline := range online {
		if userOnline.UserId == event.Data.To {
			encoded, _ := json.Marshal(event)
			log.Printf("Send message \"%s\" to user %d", encoded, event.Data.To)
			userOnline.Output.WriteMessage(websocket.TextMessage, encoded)
		}
	}
}