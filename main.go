package main

import (
	"fmt"
	"net/url"
	"strconv"
	"context"
)

var categoryID int = 57
var perPage int = 25
func GetDBKey(id, acctid int) string {
	return strconv.Itoa(id) + "|" + strconv.Itoa(acctid)

}
func GetStartRowsFromUrl(uri string) (from, to int, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	from, err = strconv.Atoi(m["StartRow"][0])
	if err != nil {
		return
	}
	count, err := strconv.Atoi(m["rowCount"][0])
	if err != nil {
		return
	}
	to = from + count - 1
	return
}
func GetProductID(uri string)(id, acctid int, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return
	}
	id, err = strconv.Atoi(m["itemid"][0])
	
	if err != nil {
		return
	}
	acctid, err = strconv.Atoi(m["acctid"][0])
	if err != nil {
		return
	}
	return
}
func main() {
	db := GetDB()
	s3Client, err := GetS3Client()
	if err != nil {
		fmt.Println("couldnt connect to s3\n", err)
		return
	}
	ctx := context.Background()
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
				continue
			}
			id, acctid, err := GetProductID(govdeals + uri)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if Exists(db, GetDBKey(id, acctid)) {
				continue
			}
			prdct, err := ParseProductPage(govdeals + uri)
			if err != nil {
				fmt.Println(err)
				continue
			}
			prdct.client = s3Client	
			prdct.ctx = ctx
			if err = SendProduct(prdct); err != nil {
				fmt.Println(err)
				continue
			}
			if err != nil {
				fmt.Println(err)
				continue
			}
			Set(db, GetDBKey(id, acctid), "sended")

		}
	}
}
