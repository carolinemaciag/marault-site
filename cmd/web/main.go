package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// ✅ Serve static files (CSS, JS, videos, images)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Routes
	mux.HandleFunc("/", homeHandler)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/home.html",
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Marault Intelligence",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
