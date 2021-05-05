package main

import (
	"flag"

	"github.com/wangyuanmike/gophercises/sitemap"
)

func main() {
	addr := flag.String("addr", "https://github.com/WangYuanMike/", "The root url of sitemap")
	depth := flag.Int("depth", 1, "Depth of sitemap")
	flag.Parse()

	sitemap.CollectSitemap(*addr, *depth)
	sitemap.GenerateXML()
}
