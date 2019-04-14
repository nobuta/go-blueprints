package main

import (
	"flag"
	"ghithub.com/nobuta/go-blueprints/chat/trace"
	"html/template"
	"log"
	"net/http"
	"os"
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
		t.tmpl = template.Must(template.ParseFiles(path))
	})
	_ = t.tmpl.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Listen AddressAndPort");
	flag.Parse()


	room := newRoom()
	room.tracer = trace.New(os.Stdout)

	http.Handle("/", &templateHandler{file:"chat.html"})
	http.Handle("/room", room)

	// channelの待ち受け開始
	go room.run()

	log.Println("Listen", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
