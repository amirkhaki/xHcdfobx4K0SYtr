package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

func render(nodes ...*html.Node) {
	for i, n := range nodes {
		var buf bytes.Buffer
		html.Render(&buf, n)
		fmt.Println("#", i)
		fmt.Println(buf.String())
	}
	fmt.Println("DONE")
}
func ParseProductPage(url string) (Product, error) {
	product := Product{}
	page, err := LoadUrl(url)
	if err != nil {
		return product, fmt.Errorf("Error during loading %s: %w", url, err)
	}
	rootNode, err := ParsePage(page)
	if err != nil {
		return product, fmt.Errorf("Error during parsing page: %w", err)
	}
	imageATag := getNodesWithAttrValue(rootNode, "a", "id", "thumb1")[0]
	for _, attr := range imageATag.Attr {
		if attr.Key == "href" {
			product.Images = append(product.Images, map[string]string{"src": "https://govdeals.com" + attr.Val})
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
	sellerTbl := getNodesWithAttrValue(rootNode, "table", "class", "table ml-1 pl-0")[0]
	tdDesc := getNodesWithAttrValue(sellerTbl, "td", "colspan", "2")[1]
	descP := getNodesWithOptions(tdDesc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && (n.Data == "P" || n.Data == "p") {
			return true
		}
		return false
	})[0]
	var buf bytes.Buffer
	err = html.Render(&buf, descP)
	if err != nil {
		return product, fmt.Errorf("Error during rendering description: %w", err)
	}
	product.Description = buf.String()
	return product, nil
}
