package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"strings"
)

var attrID string
var stopIter = false

func main() {
	url, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	attrID, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	attrID = strings.TrimSpace(attrID)
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

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) {
	if pre != nil {
		if pre(n) {
			stopIter = true
		}
	}
	for c := n.FirstChild; c != nil && stopIter; c = c.NextSibling {
		forEachNode(c, startElement, endElement)
	}
	if post != nil {
		if post(n) {
			stopIter = true
		}
	}
}

func startElement(n *html.Node) bool {
	return checkNode(n)
}

func endElement(n *html.Node) bool {
	return checkNode(n)
}

func checkNode(n *html.Node) bool {
	if findByID(n, attrID) != nil {
		fmt.Printf("%s is find\n", attrID)
		return true
	}
	return false
}

func findByID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, item := range n.Attr {
			if item.Key == "id" && item.Val == id {
				return n
			}
		}
	}
	return nil
}
