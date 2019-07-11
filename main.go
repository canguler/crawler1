package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func getLink(t html.Token, sourceURL *url.URL) (link string, ok bool) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			u, err := url.Parse(a.Val)
			if err != nil {
				log.Fatal(err)
			}
			h := u.Hostname()
			if h == "" || h == sourceURL.Hostname() {
				uu := url.URL{}
				uu.Scheme = sourceURL.Scheme
				uu.Host = sourceURL.Hostname()
				uu.Path = u.Path
				link = uu.String()
				ok = true
			}
			return
		}
	}
	return
}

func getLinks(source string) (links []string) {
	resp, _ := http.Get(source)
	b := resp.Body
	defer b.Close()
	z := html.NewTokenizer(b)

	for {
		switch z.Next() {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			t := z.Token()

			if t.Data != "a" {
				continue
			}

			link, ok := getLink(t, resp.Request.URL)
			if !ok {
				continue
			}
			links = append(links, link)
		}
	}
}

func main() {
	for _, url := range getLinks(os.Args[1]) {
		fmt.Println(url)
	}
}
