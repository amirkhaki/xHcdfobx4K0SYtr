package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendProduct(p Product) error {
	wc_prdct_endpoint := os.Getenv("WC_URL") + "/wp-json/wc/v3/products/"
	wc_prdct_endpoint += "?consumer_key=" + os.Getenv("WP_KEY") + "&consumer_secret=" + os.Getenv("WP_SECRET")
	req, err := http.NewRequest("POST", wc_prdct_endpoint, bytes.NewBuffer([]byte(p.ToJson())))
	req.SetBasicAuth(os.Getenv("WC_KEY"), os.Getenv("WC_SECRET"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error during posting product: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return nil
}
