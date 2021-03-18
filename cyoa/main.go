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

// Chapter is a chapter of the cyoa story
type Chapter struct {
	Title      string   `json:"title,omitempty"`
	Paragraphs []string `json:"story,omitempty"`
	Options    []option `json:"options,omitempty"`
}

type option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

type story map[string]Chapter

func loadStory(r io.Reader) (story, error) {
	decoder := json.NewDecoder(r)
	var s story
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

var tpl *template.Template

type handler struct {
	s      story
	t      *template.Template
	pathFn func(*http.Request) string
}

type handlerOption func(*handler)

func NewHandler(s story, opts ...handlerOption) http.Handler {
	h := handler{s, tpl, defaultPathFn}
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
		http.Error(w, "Chapter could not be executed...", http.StatusInternalServerError)
	}
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}
	return path[1:]
}

func WithTemplate(t *template.Template) handlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFn(fn func(*http.Request) string) handlerOption {
	return func(h *handler) {
		h.pathFn = fn
	}
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	if len("/story/") >= len(path) {
		return path
	}
	return path[len("/story/"):]
}

func init() {
	tpl = template.Must(template.ParseFiles("view.html"))
}

func main() {
	filename := flag.String("file", "gopher.json", "the file path of the cyoa story")
	port := flag.Int("port", 8080, "the port that HTTP server is starting at")
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
	prettyHandler := NewHandler(story, WithTemplate(prettyTpl), WithPathFn(pathFn))

	mux := http.NewServeMux()
	mux.Handle("/", defaultHandler)
	mux.Handle("/story/", prettyHandler)

	fmt.Println("Start http server at port 8080...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
