package main

import (
	"fmt"
	"net/http"
	"os"

	"strings"

	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var MaxDepth = 2

type Link struct {
	url   string
	text  string
	depth int
}

func (this *Link) String() string {
	spacer := strings.Repeat("\t", this.depth)

	return fmt.Sprintf("%s%s (%d) - %s", spacer, this.text, this.depth, this.url)
}

func (this *Link) Valid() bool {
	if this.depth >= MaxDepth {
		return false
	}

	if len(this.text) == 0 {
		return false
	}

	if len(this.url) == 0 || strings.Contains(strings.ToLower(this.url), "javascript") {
		return false
	}

	return true
}

type HttpError struct {
	original string
}

func (this HttpError) Error() string {
	return this.original
}

func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{} // 빈 slice 하나 만든다.

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()
		if token.Type == html.ErrorToken {
			break
		}

		if start != nil && token.Type == html.TextToken {
			text = fmt.Sprintf("%s%s", text, token.Data)
		}

		if token.DataAtom == atom.A { // anchor tag이면
			switch token.Type {
			case html.StartTagToken:
				if len(token.Attr) > 0 {
					start = &token
				}
			case html.EndTagToken:
				if start == nil {
					log.Warnf("Link End found without Start: %s\n", text)
					continue
				}

				link := NewLink(*start, text, depth)
				if link.Valid() {
					links = append(links, link)
					log.Debugf("Link found : %v\n", link)
				}

				start = nil
				text = ""
			}
		}
	}

	log.Debug(links)

	return links
}

func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
			//break
		}
	}

	return link
}

func downloader(url string) (resp *http.Response, err error) {
	log.Debugf("Downloading : %s\n", url)
	resp, err = http.Get(url)

	if err != nil {
		log.Debugf("Error : %s", err)
		return
	}

	if resp.StatusCode > 299 {
		err = HttpError{original: fmt.Sprintf("Error (%d): %s", resp.StatusCode, url)}
		log.Debug(err)
		return
	}

	return
}

func recurDownloader(url string, depth int) {
	page, err := downloader(url)
	if err != nil {
		log.Error(err)
		return
	}

	links := LinkReader(page, depth)

	for _, link := range links {
		fmt.Println(link)

		if depth+1 < MaxDepth {
			recurDownloader(link.url, depth+1)
		}
	}
}

func main() {
	log.SetPriorityString("info")
	log.SetPrefix("crawler")

	log.Debug(os.Args)

	if len(os.Args) < 2 {
		log.Fatalln("Missing Url arg")
	}

	recurDownloader(os.Args[1], 0)
}
