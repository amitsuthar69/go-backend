package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080", "http://localhost:4321"},
		AllowCredentials: true,
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
