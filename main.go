package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func getHref(tag html.Token) (ok bool, href string) {

	for _, attr := range tag.Attr {
		if attr.Key == "href" {
			href = attr.Val
			ok = true
		}
	}

	return

}

func crawl(url string, ch chan string, chFinished chan bool) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error failed to Crawl")
		return

	}

	defer func() {
		chFinished <- true
	}()

	body := resp.Body

	defer body.Close()

	token := html.NewTokenizer(body)

	for {
		nextTokens := token.Next()

		switch {
		case nextTokens == html.ErrorToken:
			fmt.Println("stuck with the error so returning the program")
			return
		case nextTokens == html.StartTagToken:
			tag := token.Token()

			isAnchortag := tag.Data == "a"
			if !isAnchortag {
				continue
			}

			ok, url := getHref(tag)

			if !ok {
				continue
			}

			hasProto := strings.Index(url, "https") == 0
			if hasProto {
				ch <- url
			}
		}
	}

}

func main() {

	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range seedUrls {

		go crawl(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}

	}

	fmt.Println("Found Urls", len(foundUrls), ", unique urls:")
	for url := range foundUrls {
		fmt.Println("-" + url)
	}

	close(chUrls)
}
