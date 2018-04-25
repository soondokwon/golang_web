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

var g_tplView = template.Must(template.New("view").ParseFiles("./html/base.html", "./html/view.html", "./html/header.html", "./html/footer.html"))
var g_tplEdit = template.Must(template.New("edit").ParseFiles("./html/base.html", "./html/edit.html", "./html/header.html", "./html/footer.html"))

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

	g_tplView.ExecuteTemplate(w, "base", p)
}

func edit(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, _ := load(title)

	//t, _ := template.ParseFiles("edit.html")
	//t.Execute(w, p)
	g_tplEdit.ExecuteTemplate(w, "base", p)
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
	// resource 파일 서버 위치
	http.Handle("/res/", http.StripPrefix("/res/", http.FileServer(http.Dir("res"))))

	// 일반 web handler 등록
	http.HandleFunc("/", index)
	http.HandleFunc("/view/", view)
	http.HandleFunc("/edit/", edit)
	http.HandleFunc("/save/", save)

	// http 리스너 시작
	http.ListenAndServe(":8000", nil)
}
