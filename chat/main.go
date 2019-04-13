package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
)

type templateHandler struct {
	once 	sync.Once
	file 	string
	tmpl 	*template.Template
}

func current() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		path := filepath.Join( current(), "templates", t.file)
		fmt.Println(path)
		t.tmpl = template.Must(template.ParseFiles(path))
	})
	_ = t.tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/", &templateHandler{file:"index.html"})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
