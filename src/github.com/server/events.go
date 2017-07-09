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
	case eventNames[EventNameMessageSended]:
	//assert that data interface exactly what eventMessageSended need
		if asserted, ok := data.(MessageData); ok {
			eventMessageSended(MessageSendedEvent{event, asserted})
		} else {
			log.Panic(data)
		}
	case eventNames[EventNameUserEntered]:
		eventUserEntered(UserEnteredEvent{event})
	default:
		log.Printf("Unknown event name \"%s\"", event.Name)
		log.Panic(event)
	}
}

/**
 * event occured when some user send message
 */
func eventMessageSended(event MessageSendedEvent) {
	for _, userOnline := range online {
		if userOnline.UserId == event.Data.To {
			encoded, _ := json.Marshal(event)
			log.Printf("Send message \"%s\" to user %d", encoded, event.Data.To)
			userOnline.Output.WriteMessage(websocket.TextMessage, encoded)
		}
	}
}

/**
 * event occured when some user entered
 */
func eventUserEntered(event UserEnteredEvent) {
	for _, userOnline := range online {
		if string(userOnline.UserId) != event.Source.Guid{
			encoded, _ := json.Marshal(event)
			log.Printf("Send message \"%s\" to user %s", encoded, event.Source.Guid)
			userOnline.Output.WriteMessage(websocket.TextMessage, encoded)
		}
	}
}
