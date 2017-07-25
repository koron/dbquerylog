package parser

// State represents state of the connection
type State int

const (
	None State = iota
	Handshake
	Auth
	AuthResend
	AwaitingReply
	Connected
)

type ResultState int

const (
	Fields ResultState = iota + 1
	Records
)

// Context represents the context for a connection.
type Context struct {
	State State

	ResultState ResultState
	FieldNCurr      uint64
	FieldNMax       uint64

	Data interface{}
}
