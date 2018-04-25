package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type Page struct {
	Title string
	Body  []byte
}

func (this *Page) save() error {
	f := this.Title + ".txt"

	return ioutil.WriteFile(f, this.Body, 0600)
}

func load(title string) (*Page, error) {
	f := title + ".txt"
	body, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body}, nil
}

func view(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := load(title)

	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func viewTemplete(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view-template/"):]
	p, _ := load(title)
	t, _ := template.ParseFiles("view.html")

	t.Execute(w, p)
}

func edit(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, _ := load(title)
	t, _ := template.ParseFiles("edit.html")

	t.Execute(w, p)
}

func save(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	p.save()

	// redirect : 다른 페이지로 이동
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/"):]
	if len(title) <= 0 {
		fmt.Fprintf(w, "This is index page.")
		return
	}

	p := &Page{Title: title, Body: []byte("고 언어 재미 있어요. from index page.")}
	p.save()

	fmt.Fprintf(w, "[%s.txt] generated...", title)
}

func main() {
	//p := &Page{Title: "test", Body: []byte("고 언어 재미 있어요.")}
	//p.save()

	http.HandleFunc("/", index)
	http.HandleFunc("/view/", view)
	http.HandleFunc("/view-template/", viewTemplete)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/save/", save)

	http.ListenAndServe(":8000", nil)

	// localhost:8000/view/hi
	// localhost:8000/view/index
	// localhost:8000/view/test
}
