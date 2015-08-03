package gochat

type Module interface {
	//Takes a message and evaluates whether or not the Module should act upon it
	IsValid(msg *Message, c *Channel) bool
	//Takes a message, and returns the result. If there is no result, "" should be returned
	ParseMessage(msg *Message, c *Channel) string
}
