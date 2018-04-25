package main

import (
	"context"
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
	render(w, "index", u)
}

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		msg := getFlashMsg(w, r, "message")
		if msg != "" {
			u := &User{}
			u.Errors = make(map[string]string)
			u.Errors["message"] = string(msg)
			render(w, "login", u)
			return
		}

		render(w, "login", nil)

	case "POST":
		id := r.FormValue("id")
		pw := r.FormValue("pw")

		if id == "" {
			// fmt.Fprintf(w, "id가 필요합니다.")
			setFlashMsg(w, "message", "id가 필요합니다.")
			http.Redirect(w, r, "/login", 302)
			return
		}

		if pw == "" {
			// fmt.Fprintf(w, "pw가 필요합니다.")
			setFlashMsg(w, "message", "pw가 필요합니다.")
			http.Redirect(w, r, "/login", 302)
			return
		}

		u, result := userExists(id, pw)
		if result == false {
			setFlashMsg(w, "message", "등록되지 않은 계정입니다.")
			http.Redirect(w, r, "/login", 302)
			return
		}

		setSession(u, w)
		http.Redirect(w, r, "/view", 302)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w, "session")
	http.Redirect(w, r, "/", 302)
}

func view(w http.ResponseWriter, r *http.Request) {
	id := getUserInfo(r, "id")
	pw := getUserInfo(r, "pw")
	em := getUserInfo(r, "email")
	fn := getUserInfo(r, "fname")
	ln := getUserInfo(r, "lname")

	if id == "" {
		http.Redirect(w, r, "/", 302)
		return
	}

	render(w, "view", &User{Id: id, Pw: pw, Lname: ln, Fname: fn, Email: em})
}

func userList(w http.ResponseWriter, r *http.Request) {
	users, err := makeUserList()
	if err != nil { // 오류가 발생하면...
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	render(w, "userList", users)
}

func register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		msg := getFlashMsg(w, r, "message")

		if msg != "" {
			u := &User{}
			u.Errors = make(map[string]string)
			u.Errors["message"] = string(msg)
			render(w, "register", u)
			return
		}

		u := &User{}
		render(w, "register", u)

	case "POST":
		f := r.FormValue("fName")
		l := r.FormValue("lName")
		em := r.FormValue("em")
		id := r.FormValue("id")
		pw := r.FormValue("pw")

		if id == "" {
			setFlashMsg(w, "message", "id가 필요합니다.")
			http.Redirect(w, r, "/register", 302)
			return
		}

		u := &User{Fname: f, Lname: l, Email: em, Id: id, Pw: pw}
		saveData(u)
		http.Redirect(w, r, "/login", 302)
	}
}

func render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, err := template.ParseGlob("./html/*.html")
	if err != nil { // 오류가 발생하면...
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil { // 오류가 발생하면...
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	log.Println("HTTP Server start...")

	// static file server
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", index)
	router.HandleFunc("/login", login).Methods("POST", "GET")
	router.HandleFunc("/view", view)
	router.HandleFunc("/userList", userList)
	router.HandleFunc("/register", register).Methods("POST", "GET")
	router.HandleFunc("/logout", logout).Methods("GET")

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
