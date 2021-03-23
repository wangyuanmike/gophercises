package link

import (
	"encoding/json"
	"io"
	"strings"

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

//CollectLinks() calls a recursive dfs with functional options
//It is like a practice for the functional options introduced
//in last exerciese. Here "Links" is a global variable, which
//might not be a good practice. Below the text() function is
//a common recursive function that concatenates all text from
//the subtree of a certain html node. There, no global variable
//is introduced. The common dfs recursive function should be a
//better approach, because dfs is quite concise to implement,
//the value to reuse the dfs traverse function is not so big.
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
			newLink := Link{href(n), text(n)}
			Links = append(Links, newLink)
		}
	}
}

func href(n *html.Node) string {
	var href string
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			href = attr.Val
			break
		}
	}
	return href
}

func text(n *html.Node) string {
	var ret string
	if n.Type == html.TextNode {
		ret = strings.TrimSpace(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += " " + text(c)
	}
	return strings.TrimSpace(ret)
}
