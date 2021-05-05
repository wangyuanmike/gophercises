package sitemap

import (
	"container/list"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/wangyuanmike/gophercises/link"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

var sitemap map[string]struct{}

func init() {
	sitemap = make(map[string]struct{})
}

func GenerateXML() {
	toXml := urlset{
		Xmlns: xmlns,
	}
	for page := range sitemap {
		toXml.Urls = append(toXml.Urls, loc{page})
	}
	fmt.Println()
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}

type Page struct {
	link  string
	level int
}

// BFS implementation
func CollectSitemap(addr string, depth int) {
	domain := Domain(addr)
	q := list.New()
	q.PushBack(Page{link: addr, level: 1})
	for q.Len() > 0 {
		p := q.Front()
		q.Remove(p)
		page := p.Value.(Page)
		links := CollectPageLinks(page.link)
		AddDomain(domain, &links)
		AddToSitemap(links)
		if page.level == depth {
			continue
		}
		for _, l := range links {
			level := page.level + 1
			q.PushBack(Page{link: l.Href, level: level})
		}
	}
}

func CollectPageLinks(addr string) []link.Link {
	resp, err := http.Get(addr)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	blob, err := link.CollectLinks(resp.Body)
	if err != nil {
		panic(err)
	}
	var links []link.Link
	err = json.Unmarshal(blob, &links)
	if err != nil {
		panic(err)
	}
	return links
}

func AddDomain(domain string, links *[]link.Link) {
	for _, l := range *links {
		if strings.HasPrefix(l.Href, "/") {
			l.Href = "https://" + domain + "/" + l.Href
		}
	}
}

func Domain(addr string) string {
	u, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	if u.Hostname() == "" {
		return ""
	}
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain
}

func Filter(domain string, links *[]link.Link) {
	tmp := make([]link.Link, 0)
	for _, l := range *links {
		_, ok := sitemap[l.Href]
		if ok == true {
			continue
		}
		if Domain(l.Href) == domain {
			tmp = append(tmp, l)
		}
	}
	*links = tmp
}

func AddToSitemap(links []link.Link) {
	for _, l := range links {
		sitemap[l.Href] = struct{}{}
	}
}
