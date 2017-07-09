package main

/** Defining PkgType enum */
type PkgType int
const (
	PkgTypeCommand PkgType = 0 + iota
	PkgTypeEvent
)
var pkgTypes = [...]string{
	"PkgTypeCommand",
	"PkgTypeEvent",
}
/** end define */

/** Defining SourceType enum */
type SourceType int
const (
	SourceTypeUser SourceType = 0 + iota
)
var sourceTypes = [...]string{
	"SourceTypeUser",
}
/** end define */

/** defining events type enum */
type EventName int
const (
	EventNameMessageSended EventName = 0 + iota
	EventNameUserEntered
)
var eventNames = [...]string{
	"message-sended",
	"user-entered",
}
/** end defining */

/** defining commands type enum */
type CommandName int
const (
	CommandNameSendMessage CommandName = 0 + iota
)
var commandNames = [...]string{
	"send-message",
}
/** end defining */