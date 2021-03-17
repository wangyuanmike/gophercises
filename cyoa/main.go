package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("view.html"))
}

// A Chapter is a chapter of a CYOA story
type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []option `json:"options"`
}

type option struct {
	Text    string `json:"text,omitempty"`
	Chapter string `json:"arc,omitempty"`
}

type story map[string]Chapter

func loadStory(r io.Reader) (story, error) {
	var s story
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&s); err != nil {
		return nil, err
	}
	return s, nil
}

type handler struct {
	t      *template.Template
	s      story
	pathfn func(*http.Request) string
}

func defaultPathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathfn(r)
	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

type handlerOption func(h *handler)

func WithTemplate(t *template.Template) handlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func WithPathFn(fn func(*http.Request) string) handlerOption {
	return func(h *handler) {
		h.pathfn = fn
	}
}

func pathfn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}

func NewHandler(s story, opts ...handlerOption) http.Handler {
	h := handler{tpl, s, defaultPathFn}
	for _, opt := range opts {
		opt(&h)
	}
	return h
}

func main() {
	port := flag.Int("port", 8080, "the port that http server is listening")
	filename := flag.String("file", "gopher.json", "the JSON file with the CYOA story")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(f)
	}
	defer f.Close()
	story, err := loadStory(f)
	prettyTpl := template.Must(template.ParseFiles("pretty.html"))
	handler := NewHandler(story, WithTemplate(prettyTpl), WithPathFn(pathfn))

	mux := http.NewServeMux()
	mux.Handle("/story/", handler)
	mux.Handle("/", NewHandler(story))

	fmt.Println("Start HTTP server on port 8080...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}
