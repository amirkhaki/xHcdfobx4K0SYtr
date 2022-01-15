package main

import (
	"fmt"
)

func main() {
	page, err := LoadUrl(GetCategoryUrl(57, 1, 25))
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
	for i := range elms {
		url, err := GetProductUrl(elms[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		prdct, err := ParseProductPage(govdeals + url)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(SendProduct(prdct))
	}
}
