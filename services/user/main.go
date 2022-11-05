package main

import "net/http"

func main() {
	http.ListenAndServe("0.0.0.0:8082", http.FileServer(http.Dir(".")))
}
