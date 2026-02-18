package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"
)

/* =========================
   GENERIC SERVICE PAGE HANDLER
========================= */
func servicePageHandler(templateFile string, title string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			"./internal/templates/base.html",
			"./internal/templates/"+templateFile,
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
		}{
			Title: title,
		}

		if err := tmpl.Execute(w, data); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	mux := http.NewServeMux()

	/* =========================
	   Static files
	========================= */
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	/* =========================
	   Core pages
	========================= */
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/approach", approachHandler)
	mux.HandleFunc("/executive-team", executiveTeamHandler)
	mux.HandleFunc("/contact", contactHandler)
	mux.HandleFunc("/inquire", inquireHandler)

	/* =========================
	   Services
	========================= */
	mux.HandleFunc("/services", servicesHandler)

	mux.HandleFunc(
		"/services/data-visibility-audit",
		servicePageHandler(
			"data_visibility_audit.html",
			"Data Visibility Audit | Marault Intelligence",
		),
	)

	log.Println("Starting server on :4000")
	log.Fatal(http.ListenAndServe(":4000", mux))
}

/* =========================
   HOME PAGE (ROOT ONLY)
========================= */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent "/" from catching all routes
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
   SERVICES OVERVIEW
========================= */
func servicesHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/services.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Services | Marault Intelligence",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/* =========================
   CONTACT
========================= */

func contactHandler(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles(
		"./internal/templates/base.html",
		"./internal/templates/contact.html",
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Contact | Marault Intelligence",
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
	}
}

/* =========================
   INQUIRE
========================= */

func inquireHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles(
			"./internal/templates/base.html",
			"./internal/templates/inquire.html",
		)
		if err != nil {
			log.Println(err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Title string
		}{
			Title: "Inquire | Marault Intelligence",
		}

		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodPost {

		name := r.FormValue("name")
		email := r.FormValue("email")
		company := r.FormValue("company")
		message := r.FormValue("message")
		services := r.Form["services"]

		selectedServices := strings.Join(services, ", ")

		err := sendInquiryEmail(name, email, company, selectedServices, message)
		if err != nil {
			log.Println(err)
			http.Error(w, "Unable to send message", http.StatusInternalServerError)
			return
		}

		// Render thank you page
		tmpl, _ := template.ParseFiles(
			"./internal/templates/base.html",
			"./internal/templates/thankyou.html",
		)

		data := struct {
			Title string
		}{
			Title: "Thank You | Marault Intelligence",
		}

		tmpl.Execute(w, data)
	}
}

func sendInquiryEmail(name, email, company, services, message string) error {

	from := "caroline@maraultintelligence.com"
	password := "fxhiauwzwnrqhrhk"

	to := []string{
		"caroline@maraultintelligence.com",
		"lindsey@maraultintelligence.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	body := fmt.Sprintf(
		"New Inquiry\n\nName: %s\nEmail: %s\nCompany: %s\nServices: %s\n\nMessage:\n%s",
		name, email, company, services, message,
	)

	msg := "From: " + from + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Subject: New Website Inquiry\n\n" + body

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
}





