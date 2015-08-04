package gochat

type Module interface {
	//Takes a message and evaluates whether or not the Module should act upon it
	IsValid(msg *Message, c *Channel) bool
	//Takes a message, and returns the result. If there is no result, "" should be returned
	ParseMessage(msg *Message, c *Channel) string
}

//Operation Critical Modules

//Responds PONG to PING requests
type PingResp struct {
}

func (m *PingResp) IsValid(msg *Message, c *Channel) bool {
	return msg.Cmd == "PING"
}

func (m *PingResp) ParseMessage(msg *Message, c *Channel) string {
	if len(msg.Arguments) == 2 {
		return "PONG " + msg.Arguments[1]
	} else {
		return "PONG " + msg.Arguments[0]
	}
}
