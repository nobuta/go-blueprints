package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc(
		"/",
		func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("<html>chat</html>"))
		})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
