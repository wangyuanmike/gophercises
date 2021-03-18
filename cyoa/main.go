package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type story map[string]Chapter

func loadStory(r io.Reader) (story, error) {
	var s story
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type handler struct {
	t      *template.Template
	s      story
	pathFn func(*http.Request) string
}

type handlerOption func(h *handler)

func NewHandler(s story, opts ...handlerOption) http.Handler {
	h := handler{defaultTpl, s, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFn(r)
	chapter, ok := h.s[path]
	if ok != true {
		http.Error(w, "chapter not found", http.StatusNotFound)
	}

	err := h.t.Execute(w, chapter)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Template execution failed", http.StatusInternalServerError)
	}
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func WithTemplate(t *template.Template) func(*handler) {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFn(pathFn func(r *http.Request) string) func(*handler) {
	return func(h *handler) {
		h.pathFn = pathFn
	}
}

var defaultTpl *template.Template

func init() {
	defaultTpl = template.Must(template.ParseFiles("view.html"))
}

// All code above is the complete CYOA package
// prettyPathFn() and main() could be split into another file as an execution example
func prettyPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

func main() {
	port := flag.Int("port", 8080, "the port to start the CYOA web application on")
	filename := flag.String("file", "gopher.json", "file that stores the cyoa story")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	story, err := loadStory(f)
	if err != nil {
		panic(err)
	}

	defaultHandler := NewHandler(story)
	prettyTpl := template.Must(template.ParseFiles("pretty.html"))
	prettyHandler := NewHandler(story, WithTemplate(prettyTpl), WithPathFn(prettyPathFn))

	mux := http.NewServeMux()
	mux.Handle("/", defaultHandler)
	mux.Handle("/story/", prettyHandler)

	fmt.Printf("Starting HTTP server at port %d...\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
