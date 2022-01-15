package main

import (
	"golang.org/x/net/html"
)

func GetLastPage(node *html.Node) (string, error) {
	paginationDiv := getNodesWithAttrValue(node, "div", "id", "pagination_1")[0]
	aTags := getNodesWithAttrValue(paginationDiv, "a", "style", "font-weight:bold;text-decoration:none")
	lastA := aTags[len(aTags)-1]
	for _, a := range(lastA.Attr) {
		if a.Key == "href" {
			return a.Val, nil
		}
	}
	return "", fmt.Errorf("No last Link Found")
}
