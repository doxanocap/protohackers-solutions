package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	handleFunc()
}

func handleFunc() {
	r := mux.NewRouter()
	r.HandleFunc("/", index)
	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.ListenAndServe(":1234", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("templates/index.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	temp.ExecuteTemplate(w, "index", nil)

}
