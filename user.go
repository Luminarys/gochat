package gochat

import (
	"strings"
)

type User struct {
	Nick  string
	CMode Mode
}

type Mode int

//Declare our modes, in order of precedence, i.e. a
//mode with higher value outranks one with lower value
const (
	Normal Mode = iota
	Halfop
	Operator
	Admin
	Owner
)

func parseUsers(existing map[string]*User, m string, cnick string) (others map[string]*User, me *User) {
	//This may seem redundant, but is actually quite useful
	users := make(map[string]*User)
	if existing != nil {
		users = existing
	}
	var cuser *User
	for _, nick := range strings.Split(m, " ") {
		tuser := getUser(nick)
		if tuser.Nick != cnick {
			users[tuser.Nick] = tuser
		} else {
			cuser = tuser
		}
	}
	return users, cuser
}

//Parses a nick with a mode prefix and returns a user
func getUser(s string) *User {
	if strings.HasPrefix(s, "%") {
		return &User{Nick: s[1:], CMode: Halfop}
	}
	if strings.HasPrefix(s, "@") {
		return &User{Nick: s[1:], CMode: Operator}
	}
	if strings.HasPrefix(s, "&") {
		return &User{Nick: s[1:], CMode: Admin}
	}
	if strings.HasPrefix(s, "~") {
		return &User{Nick: s[1:], CMode: Owner}
	}
	if strings.HasPrefix(s, "+") {
		return &User{Nick: s[1:], CMode: Normal}
	}
	return &User{Nick: s, CMode: Normal}
}
