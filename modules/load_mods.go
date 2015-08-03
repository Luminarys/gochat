package gcModules

import (
	"github.com/Luminarys/gochat"
)

//Loads the default provided Modules
func LoadDefaultModules(bot *gochat.Bot) {
	bot.AddModule(&PingMod{})
	um := &URLMod{}
	um.Init()
	bot.AddModule(um)
	sm := &SedMod{}
	sm.Init()
	bot.AddModule(sm)
	bot.AddModule(&QuoteMod{})
	qm := &CuteMod{}
	qm.Init()
	bot.AddModule(qm)
	bot.AddModule(&ChanJoinMod{})
}
