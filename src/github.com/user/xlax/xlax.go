package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
)

var mu sync.Mutex

type avail bool

// TODO add support for multiple rooms
var occupied avail

func (s avail) String() string {
	if s {
		return "OCCUPIED"
	} else {
		return "NOT occupied"
	}
}

type sensor int
func (s sensor) isAbout(center int) bool {
	const fuzz int = 100
	return (int(s) > (center - fuzz)) && (int(s) < (center + fuzz))
}

func room(w http.ResponseWriter, r *http.Request) {
	var room = "9w Men's"
	switch r.Method {
	case "GET":
		// Serve the resource.
		fmt.Fprintf(w, "%s is currently %s\n", room, occupied)
	case "PUT":
		// Update an existing record.
		v, err := strconv.Atoi(r.FormValue("value"))
		if err != nil {
			msg := fmt.Sprintf("'%s' is an illegal value: %s\n",
				r.FormValue("value"), err)
			http.Error(w, msg, 400)
		}
		in := sensor(v)
		// calculate new value
		var newValue avail
		if in.isAbout(200) {
			newValue = true
		} else if in.isAbout(600) {
			newValue = false
		} else {
			return
		}
		mu.Lock()
		oldValue := occupied
		occupied = newValue
		mu.Unlock()
		fmt.Fprintf(w, "%s was %s is now %s\n",
			room, oldValue, newValue)
	case "DELETE":
		// Remove the record.
		http.Error(w, "you can't remove %s: we don't have enough already\n", 405)
	default:
		// Give an error message.
		http.Error(w, "What you really want to do???\n", 405)
	}
}

func main() {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			dump,err := httputil.DumpRequestOut(r, true)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", dump)
		})
	http.HandleFunc("/room", room)
	http.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "I'm alive")
		})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
