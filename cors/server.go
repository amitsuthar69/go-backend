/* CROSS ORIGIN RESOURCE SHARING - CORS

- Cross-Origin Resource Sharing is an HTTP-header based mechanism that allows a server to indicate any origins (domain, scheme, or port)
	other than its own from which a browser should permit loading resources.
- CORS also relies on a mechanism by which browsers make a "preflight" request to the server hosting the cross-origin resource,
	in order to check that the server will permit the actual request.
- In that preflight, the browser sends headers that indicate the HTTP method and headers that will be used in the actual request.

		- For example, by using a preflight request a client might ask a server if it would allow a DELETE request, before sending a DELETE request :

		OPTIONS /resource/foo
		Access-Control-Request-Method: DELETE
		Access-Control-Request-Headers: Origin, X-Requested-With
		Origin: https://foo.bar.org

		- If the server allows it, then it will respond to the preflight request with an "Access-Control-Allow-Methods" response header, which lists DELETE:

		HTTP/1.1 204 No Content
		Connection: keep-alive
		Access-Control-Allow-Origin: https://foo.bar.org
		Access-Control-Allow-Methods: POST, GET, OPTIONS, DELETE
		Access-Control-Allow-Headers: Origin, X-Requested-With
		Access-Control-Max-Age: 86400

	- Example of a cross-origin request:
		- A JavaScript code served from https://domain-a.com uses fetch() to make a request for https://domain-b.com/data.json.
		- For security reasons, browsers restrict cross-origin HTTP requests initiated from such scripts.
		- fetch() and XMLHttpRequest follow the same-origin policy.
		- This means that a web application using those APIs can only request resources from the same origin the application was loaded from,
		  unless the response from other origins includes the right CORS headers.
*/

package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{ // [1]
			"http://localhost:8080",
			"http://localhost:4321",
		},
		AllowCredentials: true,                    // [2]
		AllowedMethods:   []string{"GET", "POST"}, // [3]
		MaxAge:           86400,                   // [4]
	})

	handler := c.Handler(mux)

	server := http.Server{
		Addr:    ":3000",
		Handler: handler,
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{ \"message\": \"Hello World\"}"))
	})

	log.Print("Server listening on :3000")
	log.Fatal(server.ListenAndServe())
}

/*

[1] : AllowedOrigins is a list of origins a cross-domain request can be executed from.
		  If the special * value is present in the list, all origins will be allowed.

[2] : AllowCredentials indicates whether the request can include user credentials like cookies,
			HTTP authentication or client side SSL certificates. The default is false.

[3] : AllowedMethods is a list of methods the client is allowed to use with cross-domain requests.
			Default value is simple methods (GET and POST).

[4] : MaxAge int: Indicates how long (in seconds) the results of a preflight request can be cached.
			The default is 0 which stands for no max age.

-------

Usage :-
❯ curl -D - -H 'Origin: http://localhost:4321' http://localhost:3000
HTTP/1.1 200 OK
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://localhost:4321
Content-Type: application/json
Vary: Origin
Date: Sat, 24 Feb 2024 16:29:59 GMT
Content-Length: 27
{ "message": "Hello World" }

-------


Can be done by using this middleware function as well :
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, r)
	})
}
handler := CorsMiddleware(mux)
*/
