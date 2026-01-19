package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// =========================
	// Static files (CSS, JS, images)
	// =========================
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// =========================
	// Routes
	// =========================
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/approach", approachHandler)
	mux.HandleFunc("/services", servicesHandler)
	mux.HandleFunc("/executive-team", executiveTeamHandler)
	mux.HandleFunc("/contact", contactHandler)

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

/* =========================
   HOME PAGE (ROOT ONLY)
========================= */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// CRITICAL: prevent "/" from catching everything
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/home.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Marault Intelligence",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/* =========================
   THE MARAULT APPROACH
========================= */
func approachHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/approach.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "The Marault Approach | Marault Intelligence",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/* =========================
   EXECUTIVE TEAM
========================= */
func executiveTeamHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/executive.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Executive Team | Marault Intelligence",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/* =========================
   PLACEHOLDERS
========================= */
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}


