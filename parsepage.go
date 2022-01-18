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
	func(){
		imageDivs := getNodesWithAttrValue(rootNode, "div", "class", "carousel-inner")
		if len(imageDivs) <1 {
			return
		}
		imageDiv := imageDivs[0]
		imgTags := getNodesWithAttrValue(imageDiv, "img", "class", "img-fluid")
		for _, img := range(imgTags) {
			for _, attr := range img.Attr {
				if attr.Key == "src" {
					product.Images = append(product.Images, map[string]string{"src": govdeals + attr.Val})
				}
			}

		}
		return
	}()
	tTdTags := getNodesWithAttrValue(rootNode, "td", "id", "asset_short_desc_id")
	if len(tTdTags) < 1 {
		return product, fmt.Errorf("No title tag")
	}
	titleTdTag := tTdTags[0]
	titleNodes := getNodesWithOptions(titleTdTag, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			return true
		}
		return false
	})
	if len(titleNodes) < 1 {
		return product, fmt.Errorf("No title node")
	}
	titleNode := titleNodes[0]
	product.Name = titleNode.Data
	bidTables := getNodesWithAttrValue(rootNode, "table", "id", "bid_tbl")
	if len(bidTables) < 1 {
		return product, fmt.Errorf("No bid Table")
	}
	bidTable := bidTables[0]
	bTags := getNodesWithOptions(bidTable, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "b" {
			return true
		}
		return false
	})
	if len(bTags) < 3 {
		return product, fmt.Errorf("no bTags")
	}
	bTag := bTags[2]
	priceNodes := getNodesWithOptions(bTag, func(n *html.Node) bool {
		if n.Type == html.TextNode {
			return true
		}
		return false
	})
	if len(priceNodes) < 1 {
		return product, fmt.Errorf("No price node")
	}
	priceNode := priceNodes[0]
	if strings.Contains(priceNode.Data, "*") {
		bTag = bTags[3]
		priceNodes = getNodesWithOptions(bTag, func(n *html.Node) bool {
			if n.Type == html.TextNode {
				return true
			}
			return false
		})
		if len(priceNodes) < 1 {
			return product, fmt.Errorf("No price node")
		}
		priceNode = priceNodes[0]

	}
	product.Price = priceNode.Data
	sellerTbls := getNodesWithAttrValue(rootNode, "table", "class", "table ml-1 pl-0")
	if len(sellerTbls) < 1 {
		return product, fmt.Errorf("no seller table")
	}
	sellerTbl := sellerTbls[0]
	tdDescs := getNodesWithAttrValue(sellerTbl, "td", "colspan", "2")
	if len(tdDescs) < 2 {
		return product, fmt.Errorf("no description td")
	}
	tdDesc := tdDescs[1]
	descPs := getNodesWithOptions(tdDesc, func(n *html.Node) bool {
		if n.Type == html.ElementNode && (n.Data == "P" || n.Data == "p") {
			return true
		}
		return false
	})
	if len(descPs) < 1 {
		return product, fmt.Errorf("no description p")
	}
	descP := descPs[0]
	descTexts := getNodesWithOptions(descP, func(n *html.Node) bool {
		return n.Type == html.TextNode
	})
	descText := ""
	for _, dP := range(descTexts) {
		descText += dP.Data
	}
	product.ShortDesc = descText
	var buf bytes.Buffer
	err = html.Render(&buf, descP)
	if err != nil {
		return product, fmt.Errorf("Error during rendering description: %w", err)
	}
	product.Description = buf.String()
	return product, nil
}
