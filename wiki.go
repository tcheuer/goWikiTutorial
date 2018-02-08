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

/*
   This is a method named save that takes its reciever p,
   a pointer to a Page.  It takes no parameters, and returns
   a value of type error.
*/
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/*
   Takes a title parameter, appends .txt to determine the filename,
   then returns a Page which contains the title and body
   and a possible error
*/
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

/*
   A simple handler for returning the phrase, "Hi there, I love %s!"
   Where %s is everything after the initial /
*/
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

/*
   Redners the page for the file <tmpl>.txt
*/
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles(tmpl + ".html")
	t.Execute(w, p)
}

/*
  Returns a text file as an article, where the title is everything after
  /view/ and the body is read from a file with the name <title>.txt
*/
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	t, err := template.ParseFiles("view.html")
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	t.Execute(w, p)
}

/*
  Handler for editing pages. Loads the page (Or creates an empty one if it doesn't exist)
  and then displays a form for creating the page.
*/
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}

	t, err := template.ParseFiles("edit.html")
	if err != nil {
		fmt.Fprintf(w, "Error:%v", err)
	} else {
		t.Execute(w, p)
	}

}

func saveHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.ListenAndServe(":8080", nil)
}
