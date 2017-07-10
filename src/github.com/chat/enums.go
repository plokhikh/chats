package main

/** Defining PkgType enum */
type PkgType int
const (
	PkgTypeCommand PkgType = 0 + iota
	PkgTypeEvent
	PkgTypeValue
)
var pkgTypes = [...]string{
	"command",
	"event",
	"value",
}
/** end define */

/** Defining SourceType enum */
type SourceType int
const (
	SourceTypeUser SourceType = 0 + iota
)
var sourceTypes = [...]string{
	"user",
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

/** Defining ValueName enum */
type ValueName int
const (
	ValueNameUsersOnline ValueName = 0 + iota
)
var valueNames = [...]string{
	"users-online",
}
/** end define */

/** defining commands type enum */
type CommandName int
const (
	CommandNameSendMessage CommandName = 0 + iota
	CommandNameGetUsersOnline
)
var commandNames = [...]string{
	"send-message",
	"get-users-online",
}
/** end defining */