package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wangyuanmike/gophercises/link"
	html "golang.org/x/net/html"
)

func traverseHTML(f io.Reader) {
	doc, err := html.Parse(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Level:\tType\tDataAtom\tData\tAttr\n")

	var fn func(*html.Node, int)
	fn = func(n *html.Node, level int) {
		//if n.Type == html.ElementNode && n.Data == "a" {
		fmt.Printf("%d:\t%d\t%v\t\t%v\t%v\n", level, n.Type, n.DataAtom, n.Data, n.Attr)
		level++
		//}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fn(c, level)
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

	fmt.Println("Call functions from package link...")
	b, err := link.CollectLinks(f)
	fmt.Printf("byte:\t%s\n", b)
}
