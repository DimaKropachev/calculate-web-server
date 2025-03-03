package handlers

import (
	"html/template"
	"net/http"
)

func FindExpression(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./server/web/find.html"))
	tmpl.Execute(w, nil)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./server/web/index.html"))
	tmpl.Execute(w, nil)
}
