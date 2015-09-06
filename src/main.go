////////////////////////////////////////////////////////////////////////////
// Porgram: Simplicity
// Purpose: Simple yet fully-functional GO web application
// Authors: Tong Sun (c) 2015, All rights reserved
////////////////////////////////////////////////////////////////////////////

/*

To build a simple Go sample application, with as less requirement as possible.

It is build on top of Goji example, artificially using all of its features.

It also depends on some other fundamental modules like toml for configuration.

*/

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/goji/param"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

////////////////////////////////////////////////////////////////////////////
// Constant and data type/structure definitions

////////////////////////////////////////////////////////////////////////////
// Global variables definitions

var cf *Config

// Note: the code below cuts a lot of corners to make the example app simple.

//::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::::
// Function main
func main() {

	cf = ConfigGet()

	//goji.Post("/logger", Logger1)
	goji.Post("/logger", Logger2)

	// Add routes to the global handler
	goji.Get("/msgs", Root)
	// Fully backwards compatible with net/http's Handlers
	//goji.Get("/", http.RedirectHandler("/msgs", 301))
	// Use your favorite HTTP verbs
	goji.Post("/msg", NewMsg)
	goji.Get("/msg", NewMsg)
	// Use Sinatra-style patterns in your URLs
	goji.Get("/users/:name", GetUser)
	// Goji also supports regular expressions with named capture groups.
	goji.Get(regexp.MustCompile(`^/m/(?P<id>\d+)$`), GetMsg)

	// Middleware can be used to inject behavior into your app. The
	// middleware for this application are defined in middleware.go, but you
	// can put them wherever you like.
	// goji.Use(PlainText)

	// If the patterns ends with "/*", the path is treated as a prefix, and
	// can be used to implement sub-routes.
	admin := web.New()
	goji.Handle("/admin/*", admin)

	// The standard SubRouter middleware helps make writing sub-routers
	// easy. Ordinarily, Goji does not manipulate the request's URL.Path,
	// meaning you'd have to repeat "/admin/" in each of the following
	// routes. This middleware allows you to cut down on the repetition by
	// eliminating the shared, already-matched prefix.
	admin.Use(middleware.SubRouter)
	// You can also easily attach extra middleware to sub-routers that are
	// not present on the parent router. This one, for instance, presents a
	// password prompt to users of the admin endpoints.
	admin.Use(SuperSecure)

	//admin.Get("/", AdminRoot)
	admin.Get("/finances", AdminFinances)
	admin.Handle("/*", http.FileServer(http.Dir(cf.Webapp.PathAdmin)))

	// Goji's routing, like Sinatra's, is exact: no effort is made to
	// normalize trailing slashes.
	goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	// Use a custom 404 handler
	goji.NotFound(NotFound)

	// Sometimes requests take a long time.
	goji.Get("/waitforit", WaitForIt)

	goji.Handle("/*", http.FileServer(http.Dir(cf.Webapp.Path)))

	// Call Serve() at the bottom of your main() function, and it'll take
	// care of everything else for you, including binding to a socket (with
	// automatic support for systemd and Einhorn) and supporting graceful
	// shutdown on SIGINT. Serve() is appropriate for both development and
	// production.
	goji.Serve()
}

////////////////////////////////////////////////////////////////////////////
// Function definitions

// Logger1 logs to console
func Logger1(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form["case"])
}

// Logger2 logs to system syslog
func Logger2(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cmd := "logger"
	for _, v := range r.Form["case"] {
		args := []string{"Processing", v}
		if err := exec.Command(cmd, args...).Run(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

// Root route (GET "/"). Print a list of msgs.
func Root(w http.ResponseWriter, r *http.Request) {
	// In the real world you'd probably use a template or something.
	io.WriteString(w, "Msger\n======\n\n")
	for i := len(Msgs) - 1; i >= 0; i-- {
		Msgs[i].Write(w)
	}
}

// NewMsg creates a new msg (POST "/msgs"). Creates a msg and redirects
// you to the created msg.
//
// To post a new msg, try this at a shell:
// $ now=$(date +'%Y-%m-%mT%H:%M:%SZ')
// $ curl -i -d "user=carl&message=Hello+World&time=$now" localhost:8000/msgs
func NewMsg(w http.ResponseWriter, r *http.Request) {
	var msg Msg

	// Parse the POST body into the Msg struct. The format is the same as
	// is emitted by (e.g.) jQuery.param.
	r.ParseForm()
	err := param.Parse(r.Form, &msg)
	msg.Time = time.Now()

	if err != nil || len(msg.Message) > cf.MaxLen {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// We make no effort to prevent races against other insertions.
	Msgs = append(Msgs, msg)
	url := fmt.Sprintf("/m/%d", len(Msgs)-1)
	http.Redirect(w, r, url, http.StatusFound)
}

// GetUser finds a given user and her msgs (GET "/user/:name")
func GetUser(c web.C, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Msger\n======\n\n")
	handle := c.URLParams["name"]
	user, ok := Users[handle]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	user.Write(w, handle)

	io.WriteString(w, "\nMsgs:\n")
	for i := len(Msgs) - 1; i >= 0; i-- {
		if Msgs[i].User == handle {
			Msgs[i].Write(w)
		}
	}
}

// GetMsg finds a particular msg by ID (GET "/msgs/\d+"). Does no bounds
// checking, so will probably panic.
func GetMsg(c web.C, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}
	// This will panic if id is too big. Try it out!
	msg := Msgs[id]

	io.WriteString(w, "Msger\n======\n\n")
	msg.Write(w)
}

// WaitForIt is a particularly slow handler (GET "/waitforit"). Try loading this
// endpoint and initiating a graceful shutdown (Ctrl-C) or Einhorn reload. The
// old server will stop accepting new connections and will attempt to kill
// outstanding idle (keep-alive) connections, but will patiently stick around
// for this endpoint to finish. How kind of it!
func WaitForIt(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is going to be legend... (wait for it)\n")
	if fl, ok := w.(http.Flusher); ok {
		fl.Flush()
	}
	time.Sleep(15 * time.Second)
	io.WriteString(w, "...dary! Legendary!\n")
}

// AdminRoot is root (GET "/admin/root"). Much secret. Very administrate. Wow.
func AdminRoot(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Msger\n======\n\nSuper secret admin page!\n")
}

// AdminFinances would answer the question 'How are we doing?'
// (GET "/admin/finances")
func AdminFinances(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Msger\n======\n\nWe're broke! :(\n")
}

// NotFound is a 404 handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
}
