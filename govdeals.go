package main

import (
	"net/http"
	"fmt"
	"io"
	"golang.org/x/net/html"
	"strings"
)
var govdealsUrl = "https://www.govdeals.com/index.cfm?fa=Main.AdvSearchResultsNew&searchPg=Category&additionalParams=true&sortOption=ad&timing=BySimple&timingType=&category=57"
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

func getNodesWithAttrValue(node *html.Node, attrKey, attrValue string) (elements []*html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			for _, a := range n.Attr {
				if a.Key == attrKye && a.Val == attrValue {
					elements = append(elements, n)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	return
}

func GetProductNodes(node *html.Node) []*html.Node {
	return getNodesWithAttrValue(node, "id" ,"boxx_row")
}
