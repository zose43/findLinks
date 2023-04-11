package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

var selectors = map[string]int{}

func main() {
	for _, url := range os.Args[1:] {
		links, err := findLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FindLinks %s\t%v", url, err)
			continue
		}
		for _, link := range links {
			println(link)
		}
		fmt.Printf("Count different selectors in doc=%d\n", len(selectors))
		for t, c := range selectors {
			fmt.Printf("%s %d\n", t, c)
		}
	}
}

func findLinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Resp code %d\t%s\t%s", resp.StatusCode, url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Invalid HTML parse %s\t%v", url, err)
	}
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}
	if n.Type == html.ElementNode {
		selectors[n.Data]++
	}
	if n.Parent != nil {
		t := strings.TrimSpace(n.Data)
		if n.Parent.Data != "script" && n.Parent.Data != "style" && n.Type == html.TextNode && t != "" {
			fmt.Printf("%s\n", n.Data)
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}
