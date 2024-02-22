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

-----

Default servemux:

The http.Handle() and http.HandleFunc() functions allow us to register routes without explicitly declaring a servemux, like this:
	http.HandleFunc("/path", pathHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))

- Behind the scenes, these functions register their routes with the default servemux.
- This is just a regular servemux like we’ve already been using, but it is initialized automatically by Go
and is stored in the http.DefaultServeMux global variable.
- When we pass nil as the second argument to http.ListenAndServe(), the server will use http.DefaultServeMux for routing.

*/

package main

import (
	"fmt"
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

func handleUserById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.PathValue("id") // [1]*
	fmt.Fprintf(w, "Hello user %s", id)
}

func handlePostCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") // [1]
		// w.WriteHeader(405)
		// w.Write([]byte("Method not allowed"))
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed) // [2]
		return
	}
	w.Write([]byte("You can create new posts here!"))
}

func user(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // [4]
	w.Write([]byte(`{"name":"Amit"}`))
}

func main() {
	mux := http.NewServeMux()

	// method 1 :
	mux.Handle("/", home{})

	// method 2 :
	mux.Handle("/user", http.HandlerFunc(user))  // [3]
	mux.HandleFunc("/user/{id}", handleUserById) // [2]*

	// method 3 :
	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) { // [3]*
		w.Write([]byte("Your posts were here..."))
	})

	mux.HandleFunc("POST /posts/create", handlePostCreate)

	server := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}
	log.Print("server listening on http://localhost:3000")
	log.Fatal(server.ListenAndServe())
}

/*
[1] : let the user know which request methods are supported for that particular URL.
			Important: Changing the response header map after a call to w.WriteHeader() or
			w.Write() will have no effect on the headers that the user receives.
			We need to make sure that your response header map contains all the headers
			we want before we call those methods.

[2] : http.Error() is a lightweight helper function which takes a given message and status code,
			then calls the w.WriteHeader() and w.Write() methods behind the scenes for us.

[3] : HandlerFunc(f) is a Handler that calls the function f.

[4] : Go will attempt to set the correct Cnotent-Type for us by content sniffing the response body
			with the http.DetectContentType() function. If this function can’t guess the content type,
			Go will fall back to setting the header Content-Type: application/octet-stream instead.

			The http.DetectContentType() function generally works quite well, but a common gotcha
			for web developers is that it can’t distinguish JSON from plain text.
			So, by default, JSON responses will be sent with a Content-Type: text/plain; charset=utf-8 header.
			You can prevent this from happening by setting the correct header manually in your handler

[x]* : These are the features introduced in Go 1.22.
			1. mux.HandleFunc("/user/{id}", handleUserById), here {id} is the wildcard entry.
			2. id := r.PathValue("id") gives the value of wildcard with name id.
			3. Method Matching, We can now explicitly mention the HTTP Method we want to allow for given patterns.
				 Any other Method except the mentioned one will return a 404 NOT FOUND
*/
