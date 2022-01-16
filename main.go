package main

import (
	"fmt"
	"net/url"
	"strconv"
)

var categoryID int = 57
var perPage int = 25

func GetStartRowsFromUrl(uri string) (from, to int, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	fmt.Println(m)
	fmt.Printf("%T\n", m["StartRow"])
	from, err = strconv.Atoi(m["StartRow"][0])
	count, err := strconv.Atoi(m["rowCount"][0])
	if err != nil {
		return
	}
	to = from + count - 1
	return
}
func main() {

	
	var urls []string
	firstUrl := GetCategoryUrl(categoryID, 1, perPage)
	urls = append(urls, firstUrl)
	page, err := LoadUrl(firstUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	node, err := ParsePage(page)
	if err != nil {
		fmt.Println(err)
		return
	}
	func() {
		lastPageUrl, err := GetLastPage(node)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, lastItem, err := GetStartRowsFromUrl(lastPageUrl)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, now, err := GetStartRowsFromUrl(firstUrl)
		if err != nil {
			return
		}
		for now < lastItem {
			newUrl := GetCategoryUrl(categoryID, now+1, perPage)
			_, now, err = GetStartRowsFromUrl(newUrl)
			if err != nil {
				continue
			}
			urls = append(urls, newUrl)
		}

	}()
	fmt.Println(urls)
	for _, u := range urls {
		fmt.Println(u)
		page, err = LoadUrl(u)
		if err != nil {
			fmt.Println(err)
			continue
		}
		node, err = ParsePage(page)
		if err != nil {
			fmt.Println(err)
			continue
		}
		elms := GetProductNodes(node)
		for i := range elms {
			uri, err := GetProductUrl(elms[i])
			if err != nil {
				fmt.Println(err)
				return
			}
			prdct, err := ParseProductPage(govdeals + uri)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println(SendProduct(prdct))
		}
	}
}
