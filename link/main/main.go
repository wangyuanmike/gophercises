package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wangyuanmike/gophercises/link"
	html "golang.org/x/net/html"
)

func htmlTree(f io.Reader) {
	doc, err := html.Parse(f)
	if err != nil {
		panic(err)
	}

	var fn func(*html.Node, int)
	fn = func(n *html.Node, padding int) {
		for i := 0; i < padding; i++ {
			fmt.Printf("\t")
		}
		msg := n.Data
		if n.Type == html.ElementNode {
			msg = fmt.Sprintf("<%v>", msg)
		}
		fmt.Println(msg)

		padding++
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c, padding)
		}
	}
	fn(doc, 0)
}

func main() {
	filename := flag.String("file", "../ex1.html", "The html file to be parsed")
	flag.Parse()

	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	//htmlTree(f)
	//
	//The same html file could not be parsed by html.Parse() twice,
	//otherwise the second parse would only return the top 3 level
	//of the html tree structure from the root node.
	//Therefore the htmlTree(f) must be commented out, although
	//it is just to print out the html tree structure fore pre-check.
	//This might be a bug of the html x package

	b, err := link.CollectLinks(f)
	fmt.Printf("%s\n", b)
}
