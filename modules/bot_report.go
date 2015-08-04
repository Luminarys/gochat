package gcModules

import (
	"github.com/Luminarys/gochat"
)

//Bot report - Bot will respond when .bots is used
type BotReportMod struct {
}

func (m *BotReportMod) IsValid(msg *gochat.Message, c *gochat.Channel) bool {
	return msg.Text == ".bots"
}

func (m *BotReportMod) ParseMessage(msg *gochat.Message, c *gochat.Channel) string {
	return "Reporting in! [Go]"
}
