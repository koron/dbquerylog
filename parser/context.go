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

// Context represents the context for a connection.
type Context struct {
	State State
}
