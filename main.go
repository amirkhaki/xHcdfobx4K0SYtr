package main

import (
	"fmt"
)

func main() {
	page, err := LoadUrl(govdealsUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := ParsePage(page)
	if err != nil {
		fmt.Println(err)
		return
	}
	elms := GetProductNodes(node)
	for i := range(elms) {
		url, err := GetProductUrl(elms[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		prdct, err := ParseProductPage(govdeals+url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(prdct)
	}
}
