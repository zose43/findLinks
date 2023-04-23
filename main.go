package main

import (
	"fmt"
	"github.com/zose43/findLinks/links"
	"os"
)

func main() {
	breadthFirst(crawl, os.Args[1:])
}

func breadthFirst(f func(item string) []string, worklist []string) {
	seem := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seem[item] {
				seem[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	println(url)
	items, err := links.Extract(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FindLinks %s\t%v", url, err)
	}
	return items
}
