package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
)

var depth int

func main() {
	url, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	err := findLinks(url[:len(url)-1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "FindLinks %s\t%v", url, err)
	}
}

func findLinks(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Resp code %d\t%s\t%s", resp.StatusCode, url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return fmt.Errorf("Invalid HTML parse %s\t%v", url, err)
	}
	forEachNode(doc, startElement, endElement)
	return nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, startElement, endElement)
	}
	if post != nil {
		post(n)
	}
}

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s>\n", depth*2, "", n.Data)
		depth++
	}
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
	}
}
