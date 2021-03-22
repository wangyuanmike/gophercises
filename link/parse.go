package link

import (
	"encoding/json"
	"io"

	html "golang.org/x/net/html"
)

// Encode the link map in json. Here is an example:
/*
[
	{"Href": "/dog", "Text": "Dog's homepage"},
	{"Href": "/cat", "Text": "Cat's homepage"},
	{"Href": "/horse", "Text": "Horse's homepage"}
]
*/

type Link struct {
	Href string `json:"Href,omitempty"`
	Text string `json:"Text,omitempty"`
}

var Links []Link

func CollectLinks(f io.Reader) ([]byte, error) {
	doc, err := html.Parse(f)
	if err != nil {
		return nil, err
	}

	dfs(doc, withCollectLink())

	b, err := json.Marshal(Links)
	return b, err
}

func dfs(n *html.Node, opts ...dfsOption) {
	for _, opt := range opts {
		opt(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		dfs(c, opts...)
	}
}

type dfsOption func(*html.Node)

func withCollectLink() func(*html.Node) {
	return func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			newLink := Link{n.Attr[0].Val, n.FirstChild.Data}
			Links = append(Links, newLink)
		}
	}
}
