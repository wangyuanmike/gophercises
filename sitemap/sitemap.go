package sitemap

import (
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

var sitemap map[string]interface{}

func init() {
	sitemap = make(map[string]interface{})
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

func CollectSitemap(addr string, depth int, count int) {
	if count == depth {
		return
	}

	domain := Domain(addr)
	links := CollectPageLinks(addr)
	AddDomain(domain, &links)
	Filter(domain, &links)
	if len(links) == 0 {
		return
	}
	AddToSitemap(links)
	count++

	for _, link := range links {
		CollectSitemap(link.Href, depth, count)
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
	for _, link := range *links {
		if Domain(link.Href) == "" {
			link.Href = "https://" + domain + "/" + link.Href
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
	for _, link := range *links {
		_, ok := sitemap[link.Href]
		if ok == true {
			continue
		}
		if Domain(link.Href) == domain {
			tmp = append(tmp, link)
		}
	}
	*links = tmp
}

func AddToSitemap(links []link.Link) {
	for _, link := range links {
		sitemap[link.Href] = ""
	}
}
