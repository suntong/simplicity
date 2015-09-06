package main

import (
	"fmt"
	"io"
	"time"
)

// A Msg is a 140-character micro-blogpost that has no resemblance whatsoever
// to the noise a bird makes.
type Msg struct {
	User    string    `param:"user"`
	Message string    `param:"message"`
	Time    time.Time `param:"time"`
}

// Store all our msgs in a big list in memory, because, let's be honest, who's
// actually going to use a service that only allows you to post 140-character
// messages?
var Msgs = []Msg{
	{"carl", "Welcome to Msger!", time.Now()},
	{"alice", "Wanna know a secret?", time.Now()},
	{"bob", "Okay!", time.Now()},
	{"eve", "I'm listening...", time.Now()},
}

// Write out a representation of the msg
func (g Msg) Write(w io.Writer) {
	fmt.Fprintf(w, "%s at %s\n  %s\n---\n", g.User,
		g.Time.Format(time.UnixDate), g.Message)
}

// A User is a person. It may even be someone you know. Or a rabbit. Hard to say
// from here.
type User struct {
	Name, Bio string
}

// All the users we know about! There aren't very many...
var Users = map[string]User{
	"alice": {"Alice in Wonderland", "Eating mushrooms"},
	"bob":   {"Bob the Builder", "Making children dumber"},
	"carl":  {"Carl Jackson", "Duct tape aficionado"},
}

// Write out the user
func (u User) Write(w io.Writer, handle string) {
	fmt.Fprintf(w, "%s (@%s)\n%s\n", u.Name, handle, u.Bio)
}
