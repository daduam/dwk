package main

import (
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
<!DOCTYPE html>
<html>
<head><title>Todos</title></head>
<body>
<h1>Todos</h1>
<p>Welcome to the Todos app!</p>
</body>
</html>`))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", indexHandler)

	log.Printf("Server started in port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
