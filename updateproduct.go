package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

type UpProduct struct {
	ID          int             `json:"-"`
	Date        time.Time       `json:"-"`
	Name        string          `json:"name"`
	Status      string          `json:"status"`
	Description string          `json:"description"`
	ShortDesc   string          `json:"short_description"`
	Price       string          `json:"regular_price"`
	ctx         context.Context `json:"-"`
	Categories  []Category      `json:"categories"`
}

func (p UpProduct) GetPrice() (price string) {
	price = strings.Replace(p.Price, "$", "", -1)
	price = strings.Replace(price, ",", "", -1)
	return
}
func (p UpProduct) ToJson() string {
	p.Status = "publish"
	if os.Getenv("WC_DEVELOPING") != "" {
		p.Status = "pending"
	}
	if p.Date.Before(time.Now()) {
		p.Status = "private"
	}

	p.Price = p.GetPrice()
	jsonStr, _ := json.Marshal(p)
	return string(jsonStr)
}

func (p UpProduct) ToValues() (values url.Values) {
	values = make(url.Values)
	values.Set("name", p.Name)
	values.Set("status", "pending")
	values.Set("sale_price", p.Price)
	values.Set("description", p.Description)
	return
}
func (p UpProduct) String() string {
	return fmt.Sprintf("name:%s\ndescription:%s\nprice:%s\n", p.Name, p.Description, p.Price)
}

func mapProductToUpProduct(p Product) (up UpProduct) {
	up.ID = p.ID
	up.Date = p.Date
	up.Name = p.Name
	up.Status = p.Status
	up.Description = p.Description
	up.ShortDesc = p.ShortDesc
	up.Price = p.Price
	up.ctx = p.ctx
	up.Categories = p.Categories

	return
}
