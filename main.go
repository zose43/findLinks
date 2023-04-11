package main

import (
	"fmt"
	"golang.org/x/net/html"
	"os"
	"strings"
)

var selectors = map[string]int{}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FindLinks: %v\n", err)
		os.Exit(1)
	}
	for _, link := range visit(nil, doc) {
		println(link)
	}
	fmt.Printf("Count different selectors in doc=%d\n", len(selectors))
	for t, c := range selectors {
		fmt.Printf("%s %d\n", t, c)
	}
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
