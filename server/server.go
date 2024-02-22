/*
ServeMux is an HTTP request multiplexer / A Router.
It’s responsible for matching the URL in the request to an appropriate handler and executing it

we attach an URL and a handler to a ServeMux instance using the "Handle" and "HandleFunc" methods.

1. Handle :- func (mux *ServeMux) Handle(pattern string, handler Handler)
- It accepts a String (URL) and an "http.Handler".
- An "http.Handler" is an interface with "serveHTTP" method.
	type Handler interface {
	    ServeHTTP(ResponseWriter, *Request)
	}
- In order to use the Handle method, we can create a handler as a struct and implement a "ServeHTTP" method on it.
- Any struct that defines the serveHTTP() method on it, implements the Handler interface and hence becomes a http.Handler
- Example: the home{} struct in below code.

2. HandleFunc :- func (mux *ServeMux) HandleFunc(pattern string, handler func(ResponseWriter, *Request))
- Unlike the Handle method, HandleFunc accepts the handler implementation in the form of a function.
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("hello"))
	})

Now we have a multiplexer which can respond if a user navigates to the different routes of our service.
We can now tie it all together with a Server.
- We do this by creating a new instance of a Server:
	server := http.Server{Addr: "host:port", Handler: http.Handler}

-------

- Go’s servemux supports two different types of URL patterns: fixed paths and subtree paths.
- Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.

1. Our two patterns below — "/user" and "/posts" — are both examples of fixed paths.
Fixed path patterns like these are only matched (and the corresponding handler is called)
when the request URL path exactly matches the fixed path.

2. In contrast, our pattern "/" is an example of a subtree path (because it ends in a trailing slash).
Subtree path patterns are matched (and the corresponding handler is called)
whenever the start of a request URL path matches the subtree path.
The pattern "/" essentially means match a single slash, followed by anything (or nothing at all).

It’s not possible to change the behavior of Go’s servemux to do this,
but we can include a simple check for the "/" in the hander.
*/

package main

import (
	"log"
	"net/http"
)

type home struct{}

func (h home) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The "/" pattern matches everything, so we need to check that we're at the root here.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("2024 is for Go!"))
}

func user(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello User"))
}

func main() {
	mux := http.NewServeMux()

	// method 1 :
	mux.Handle("/", home{})

	// method 2 :
	mux.Handle("/user", http.HandlerFunc(user)) // HandlerFunc(f) is a Handler that calls the function f

	// method 3 :
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Your posts were here..."))
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	log.Print("server listening on http://localhost:3000")
	log.Fatal(server.ListenAndServe())
}
