package main

import (
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func ParseProductPage(url string) (Product, error) {
	var product = Product{}
	page, err := LoadUrl(url)
	if err != nil {
		return product, fmt.Errorf("Error during loading %s: %w", url, err)
	}
	rootNode, err := ParsePage(page)
	if err != nil {
		return product, fmt.Errorf("Error during parsing page: %w", err)
	}
	imageATag := getNodesWithAttrValue(rootNode, "a", "id", "thumb1")[0]
	for _, attr := range(imageATag.Attr) {
		if attr.Key == "href" {
			product.Image = "https://govdeals.com"+attr.Val
		}
	}
	titleTdTag := getNodesWithAttrValue(rootNode, "td", "id", "asset_short_desc_id")[0]
	titleNode := getNodesWithOptions(titleTdTag, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			return true
		}
		return false
	})[0]
	product.Name = titleNode.Data
	bidTable := getNodesWithAttrValue(rootNode, "table", "id", "bid_tbl")[0]
	bTags := getNodesWithOptions(bidTable, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "b" {
			return true
		}
		return false
	})
	bTag := bTags[2]
	priceNode := getNodesWithOptions(bTag, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			return true
		}
		return false
	})[0]
	if strings.Contains(priceNode.Data, "*") {
		bTag = bTags[3]
		priceNode = getNodesWithOptions(bTag, func(n *html.Node) bool {
			if n.Type == html.TextNode {
				return true
			}
			return false
		})[0]

	}
	product.Price = priceNode.Data
	return product, nil 
}
