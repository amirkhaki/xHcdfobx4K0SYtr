package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/minio/minio-go/v7"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"time"
	"os"
	"strings"
	"log"
)

var govdeals = "https://www.govdeals.com/"

func GetCategoryUrl(category, start, count int) string {
	return fmt.Sprintf("%sindex.cfm?fa=Main.AdvSearchResultsNew&searchPg=Category&additionalParams=true&sortOption=ad&timing=BySimple&timingType=&category=%d&rowCount=%d&StartRow=%d", govdeals, category, count, start)
}
type Category struct{
	ID int `json:"id"`
}
type Product struct {
	ID          int                 `json:"-"`
	Date time.Time `json:"-"`
	Name        string              `json:"name"`
	Status      string              `json:"status"`
	Description string              `json:"description"`
	ShortDesc   string              `json:"short_description"`
	Price       string              `json:"regular_price"`
	Images      []map[string]string `json:"images"`
	client      *minio.Client       `json:"-"`
	ctx         context.Context     `json:"-"`
	Categories []Category `json:"categories"`
}

func (p Product) UploadImages() []map[string]string {
	var images []map[string]string
	if len(p.Images) != 0 {
		var tmp []map[string]string
		tmp = append(tmp, p.Images[0])
		p.Images = tmp
	}
	for _, img := range p.Images {
		if src, ok := img["src"]; ok {
			f, err := DownloadFile(src, "jpg")
			defer os.Remove(f.Name())
			if err != nil {
				log.Println(err)
				continue
			}
			fileName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
			err = UploadToS3(p.client, p.ctx, f.Name(), fileName, "image/jpeg")
			if err != nil {
				log.Println(err)
				continue
			}
			images = append(images, map[string]string{"src": GetObjectUrl(fileName)})
		}
	}
	return images
}

func (p Product) DeleteImages() ( err error ) {
	for _, img := range p.Images {
		if src, ok := img["src"]; ok {
			srcParts := strings.Split(src,  "/")
			fileName := srcParts[len(srcParts) - 1]
			if er := DeleteFromS3(p.client, p.ctx, fileName); er != nil {
				err = fmt.Errorf("DeleteImage err: %s: %w", er.Error(), err)
			}
		}
	}
	return
}

func (p Product) GetPrice() (price string) {
	price = strings.Replace(p.Price, "$", "", -1)
	price = strings.Replace(price, ",", "", -1)
	return
}
func (p Product) ToJson() string {
	p.Status = "publish"
	if os.Getenv("WC_DEVELOPING") != "" {
		p.Status = "pending"
		p.Images = []map[string]string{}
	} else {
		p.Images = p.UploadImages()
	}
	if p.Date.Before(time.Now()) {
		p.Status = "private" 
	}

	p.Price = p.GetPrice()
	jsonStr, _ := json.Marshal(p)
	return string(jsonStr)
}
func (p Product) ToValues() (values url.Values) {
	values = make(url.Values)
	values.Set("name", p.Name)
	values.Set("status", "pending")
	values.Set("sale_price", p.Price)
	values.Set("description", p.Description)
	return
}
func (p Product) String() string {
	return fmt.Sprintf("name:%s\ndescription:%s\nprice:%s\nimage:%s\n", p.Name, p.Description, p.Price, p.Images)
}

type nodeOption func(*html.Node) bool

func LoadUrl(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error during creating request to %s: %w", url, err)
	}
	req.Header.Set("Host", "www.govdeals.com")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 10; M2006C3LG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,fa;q=0.8")

	resp, err := http.DefaultClient.Do(req)
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
		for _, op := range options {
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

func getNodesWithAttrValue(node *html.Node, elm, attrKey, attrValue string) (elements []*html.Node) {
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
	return getNodesWithAttrValue(node, "div", "id", "boxx_row")
}

func GetProductUrl(product *html.Node) (string, error) {
	aNode := getNodesWithAttrValue(product, "div", "id", "result_col_2")[0]
	aTag := getNodesWithTagName(aNode, "a")[0]
	for _, attr := range aTag.Attr {
		if attr.Key == "href" {
			return attr.Val, nil
		}
	}
	return "", fmt.Errorf("No product url")

}
