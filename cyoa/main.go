package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

// An Arc is a complete section of CYOA book
type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

// Book is the top struct representing a CYOA book
type Book struct {
	Intro     Arc `json:"intro,omitempty"`
	NewYork   Arc `json:"new-york,omitempty"`
	Debate    Arc `json:"debate,omitempty"`
	SeanKelly Arc `json:"sean-kelly,omitempty"`
	MarkBates Arc `json:"mark-bates,omitempty"`
	Denver    Arc `json:"denver,omitempty"`
	Home      Arc `json:"home,omitempty"`
}

func exit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func loadJSON(fileName string) []byte {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		exit(err)
	}
	defer jsonFile.Close()

	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		exit(err)
	}
	return jsonByte
}

func parseJSON(jsonByte []byte) Book {
	var book Book
	err := json.Unmarshal(jsonByte, &book)
	if err != nil {
		exit(err)
	}
	return book
}

func getArc(book Book, path string) Arc {
	arc := Arc{}

	switch path {
	case "/intro/":
		arc = book.Intro
	case "/new-york/":
		arc = book.NewYork
	case "/debate/":
		arc = book.Debate
	case "/sean-kelly/":
		arc = book.SeanKelly
	case "/mark-bates/":
		arc = book.MarkBates
	case "/denver/":
		arc = book.Denver
	case "/home/":
		arc = book.Home
	}
	return arc
}

// httpCache is to cache init information before starting http server
type httpCache struct {
	Template *template.Template
	book     Book
}

func (c *httpCache) parseTemplate(templateFile string) {
	c.Template = template.Must(template.ParseFiles(templateFile))
}

func (c *httpCache) loadBook(bookFile string) {
	c.book = parseJSON(loadJSON(bookFile))
}

func viewHandler(c httpCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		err := c.Template.ExecuteTemplate(w, c.Template.Name(), getArc(c.book, r.URL.Path))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	var c httpCache
	c.parseTemplate("view.html")
	c.loadBook("gopher.json")

	http.HandleFunc("/intro/", viewHandler(c))
	http.HandleFunc("/new-york/", viewHandler(c))
	http.HandleFunc("/debate/", viewHandler(c))
	http.HandleFunc("/sean-kelly/", viewHandler(c))
	http.HandleFunc("/mark-bates/", viewHandler(c))
	http.HandleFunc("/denver/", viewHandler(c))
	http.HandleFunc("/home/", viewHandler(c))

	fmt.Println("Start listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}
