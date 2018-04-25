package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var (
	cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	router        = mux.NewRouter()
)

func index(w http.ResponseWriter, r *http.Request) {

	u := &User{}
	tmpl, _ := template.ParseFiles("./html/index.html", "./html/header.html", "./html/navbar.html", "./html/footer.html")
	err := tmpl.ExecuteTemplate(w, "index", u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &User{}
		tmpl, _ := template.ParseFiles("./html/login.html", "./html/header.html", "./html/navbar.html", "./html/footer.html")
		tmpl.ExecuteTemplate(w, "login", u)

	case "POST":
		id := r.FormValue("id")
		pw := r.FormValue("pw")

		redirect := "/"
		if id == getUserInfo(r, "id") &&
			pw == getUserInfo(r, "pw") {
			redirect = "/view"
			http.Redirect(w, r, redirect, 302)
		} else {
			fmt.Fprintf(w, "Login failed...")
		}
	}
}

func unRegister(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func view(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./html/view.html", "./html/header.html", "./html/navbar.html", "./html/footer.html")
	id := getUserInfo(r, "id")
	pw := getUserInfo(r, "pw")
	em := getUserInfo(r, "email")
	fn := getUserInfo(r, "fname")
	ln := getUserInfo(r, "lname")

	if id != "" {
		err := tmpl.ExecuteTemplate(w, "view", &User{Id: id, Pw: pw, Lname: ln, Fname: fn, Email: em})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u := &User{}
		tmpl, _ := template.ParseFiles("./html/register.html", "./html/header.html", "./html/navbar.html", "./html/footer.html")
		tmpl.ExecuteTemplate(w, "register", u)

	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("em")
		id := r.FormValue("id")
		pw := r.FormValue("pw")

		u := &User{Fname: f, Lname: l, Email: em, Id: id, Pw: pw}
		setSession(u, w)
		http.Redirect(w, r, "/view", 302)
	}
}

func main() {
	log.Println("HTTP Server start...")

	// static file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", index)
	router.HandleFunc("/login", login).Methods("POST", "GET")
	router.HandleFunc("/view", view)
	router.HandleFunc("/register", register).Methods("POST", "GET")
	router.HandleFunc("/unregister", unRegister).Methods("GET")

	srv := &http.Server{
		Addr: "0.0.0.0:8000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("HTTP Server off...")

	os.Exit(0)
}
