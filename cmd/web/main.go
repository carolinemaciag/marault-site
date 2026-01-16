package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// =========================
	// Static files (CSS, JS, images, videos)
	// =========================
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// =========================
	// Routes
	// =========================
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/approach", approachHandler)
	mux.HandleFunc("/services", servicesHandler)
	mux.HandleFunc("/team", teamHandler)
	mux.HandleFunc("/contact", contactHandler)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// =========================
// HOME PAGE
// =========================
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

// =========================
// THE MARAULT APPROACH
// =========================
func approachHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/approach.html",
	)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := struct {
		Title string
	}{
		Title: "The Marault Approach | Marault Intelligence",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}

