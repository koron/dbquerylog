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
	PrepareParamsAndColumns
	PrepareParams
	PrepareColumns
)

type CommandType int

const (
	Sleep       CommandType = 0x00
	Quit                    = 0x01
	InitDB                  = 0x02
	Query                   = 0x03
	FieldList               = 0x04
	CreateDB                = 0x05
	DropDB                  = 0x06
	Refresh                 = 0x07
	Shutdown                = 0x08
	Statistics              = 0x09
	ProcessInfo             = 0x0a

	Prepare      = 0x16
	Execute      = 0x17
	SendLongData = 0x18
	Close        = 0x19
	Reset        = 0x1a
	SetOption    = 0x1b
	Fetch        = 0x1c
)

type Command interface {
	CommandType() CommandType
}

// Context represents the context for a connection.
type Context struct {
	ClientFlags ClientFlags

	// Compression support
	WillCompress bool
	Compressing  bool

	// (client) Query attributes enabled
	QueryAttributes bool

	State State

	// Server status
	ResultState ResultState
	FieldNCurr  uint64
	FieldNMax   uint64

	// Prepare statements
	PreparedStmts map[uint32]Stmt

	// Client status
	LastCommand CommandType

	Data interface{}
}

func newContext() *Context {
	return &Context{
		PreparedStmts: map[uint32]Stmt{},
	}
}

func (ctx *Context) IsClientDeprecateEOF() bool {
	return ctx.ClientFlags&ClientDeprecateEOF != 0
}

func (ctx *Context) addStmt(st Stmt) {
	ctx.PreparedStmts[st.ID] = st
}

type Stmt struct {
	ID         uint32
	NumParams  uint16
	NumColumns uint16
}
