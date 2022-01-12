package main

import (
	"net/http"
	"fmt"
	"io"
	"golang.org/x/net/html"
	"strings"
)
var govdealsUrl = "https://www.govdeals.com/index.cfm?fa=Main.AdvSearchResultsNew&searchPg=Category&additionalParams=true&sortOption=ad&timing=BySimple&timingType=&category=57"

var govdeals = "https://www.govdeals.com/"

type Product struct {
	Name string
	Slug string
	DateCreated string
	Description string
	RegularPrice string 
	Image string

}

type nodeOption func(*html.Node) bool
func LoadUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error during getting %s: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error during reading response body: %w", err)
	}
	return string(body), nil
}

// TODO: select wnated parts with regex and only parse them
func ParsePage(page string) (*html.Node, error) {
	node, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return nil, fmt.Errorf("Error during parsing page: %w", err)
	}
	return node, nil
}

func getNodesWithOptions(node *html.Node, options ...nodeOption) (elements []*html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		doAppend := true
		for _, op := range(options) {
			if !op(n) {
				doAppend = false
				break
			}
		}
		if doAppend {
			elements = append(elements, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return
}

func getNodesWithAttrValue(node *html.Node,elm, attrKey, attrValue string) (elements []*html.Node) {
	var f nodeOption
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == elm {
			for _, a := range n.Attr {
				if a.Key == attrKey && a.Val == attrValue {
					return true
				}
			}
		}
		return false
	}
	return getNodesWithOptions(node, f)
	
}

func getNodesWithTagName(node *html.Node, tag string) (elements []*html.Node) {
	var f nodeOption
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == tag {
			return true
		}
		return false
	}
	return getNodesWithOptions(node, f)
	
}

func GetProductNodes(node *html.Node) []*html.Node {
	return getNodesWithAttrValue(node, "div" ,"id" ,"boxx_row")
}


func GetProductUrl(product *html.Node) (string, error) {
	aNode := getNodesWithAttrValue(product, "div", "id", "result_col_2")[0]
	aTag := getNodesWithTagName(aNode, "a")[0]
	for _, attr := range(aTag.Attr) {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("No product url")


}
