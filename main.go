package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"flag"
	"log"
	"os"
)

var categoryID int = 67
var distinationCategoryID int = 180
var perPage int = 25

func init() {
	log.SetOutput(os.Stdout)
	flag.IntVar(&categoryID, "category", 67, "category id in govdeals")
	flag.IntVar(&distinationCategoryID, "dist", 180, "distination category id in your wordpress websitw")
	flag.Parse()
	debug(categoryID, distinationCategoryID)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}


func unmarshal(jsn string) (map[string]interface{}, error) {
	var rtrn map[string]interface{}
	if err := json.Unmarshal([]byte(jsn), &rtrn); err != nil {
		return nil, fmt.Errorf("coudn't unmarshal: %w", err)
	}
	return rtrn, nil
}
func GetDBKey(id, acctid int) string {
	return strconv.Itoa(id) + "-" + strconv.Itoa(acctid)

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
func GetProductID(uri string) (id, acctid int, err error) {
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
	sended := []string{}
	s3Client, err := GetS3Client()
	if err != nil {
		log.Println("couldnt connect to s3\n", err)
		return
	}
	ctx := context.Background()
	var urls []string
	firstUrl := GetCategoryUrl(categoryID, 1, perPage)
	urls = append(urls, firstUrl)
	page, err := LoadUrl(firstUrl)
	if err != nil {
		log.Println(err)
		return
	}
	node, err := ParsePage(page)
	if err != nil {
		log.Println(err)
		return
	}
	func() {
		lastPageUrl, err := GetLastPage(node)
		if err != nil {
			log.Println(err)
			return
		}
		_, lastItem, err := GetStartRowsFromUrl(lastPageUrl)
		if err != nil {
			log.Println(err)
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
		page, err = LoadUrl(u)
		if err != nil {
			log.Println(err)
			continue
		}
		node, err = ParsePage(page)
		if err != nil {
			log.Println(err)
			continue
		}
		elms := GetProductNodes(node)
		for i := range elms {
			uri, err := GetProductUrl(elms[i])
			if err != nil {
				log.Println(err)
				continue
			}
			id, acctid, err := GetProductID(govdeals + uri)
			if err != nil {
				log.Println(err)
				continue
			}
			dbKey := GetDBKey(id, acctid)
			prdct, err := ParseProductPage(govdeals + uri)
			if err != nil {
				log.Println(err)
				continue
			}
			prdct.client = s3Client
			prdct.ctx = ctx
			if Exists(db, dbKey) {
				pidStr, err := Get(db, dbKey)
				if err != nil {
					log.Println(err)
					continue
				}
				pid, err := strconv.Atoi(pidStr)
				if err != nil {
					log.Println(err)
					continue
				}
				prdct.ID = pid
				_, err = UpdateProduct(prdct)
				if err != nil {
					log.Println(err)
					continue
				}
				sended = append(sended, dbKey)
				log.Println("product updated: ", prdct.ID)
				err = prdct.DeleteImages()
				if err != nil {
					log.Println(err)
				}
				continue
			}
			prdct.Categories = append(prdct.Categories, Category{ID: distinationCategoryID})
			resp, err := SendProduct(prdct)
			if err != nil {
				log.Println(err)
				continue
			}
			respMap, err := unmarshal(resp)
			if err != nil {
				log.Println(err)
				continue
			}
			if id, ok := respMap["id"].(float64); ok {
				idStr := strconv.Itoa(int(id))
				log.Println("product sended: ", idStr)
				err = Set(db, dbKey, idStr)
				if err != nil {
					log.Println(err)
				}
				sended = append(sended, dbKey)
				err = prdct.DeleteImages()
				if err != nil {
					log.Println(err)
				}
				continue

			} else {
				debug(respMap)
			}

		}
	}
	deleted := []string{}
	for k := range(db.Keys(nil)) {
		if !contains(sended, k) {
			dpidStr, err := Get(db, k)
			if err != nil {
				log.Println(err)
				continue
			}
			dpid, err := strconv.Atoi(dpidStr)
			if err != nil {
				log.Println(err)
				continue
			}
			if _, err = DeleteProduct(dpid); err != nil {
				log.Println(err)
				continue
			}
			Delete(db, k)
			deleted = append(deleted, k)
			log.Println("product deleted: ", dpid)
			continue

		}
	}
	log.Println("products sended: ", sended)
	log.Println("products deleted: ", deleted)
}
