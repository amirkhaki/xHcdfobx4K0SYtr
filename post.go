package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func validateProduct(p Product) error {
	keywords := os.Getenv("X_KEYWORDS")
	if keywords == "" {
		return nil
	}
	// kyewords slice
	KWs := strings.Split(keywords, ",")
	if contains := func() bool {
		name := strings.ToLower(p.Name)
		for _, kw := range KWs {
			kw = strings.ToLower(kw)
			if strings.Contains(name, kw) {
				debug("found\nname in lower case is ", name)
				debug("kw is ", kw)
				return true
			}
		}
		return false
	}(); !contains {
		return fmt.Errorf("product name doesnt contain any of keywords")
	}
	return nil

}

func DeleteProduct(pid int) (string, error){
	wc_up_prdct_endpoint := os.Getenv("WC_URL") + "/wp-json/wc/v3/products/"
	wc_up_prdct_endpoint += fmt.Sprintf("%d", pid)
	wc_up_prdct_endpoint += "?consumer_key=" + os.Getenv("WP_KEY") + "&consumer_secret=" + os.Getenv("WP_SECRET")
	req, err := http.NewRequest("DELETE", wc_up_prdct_endpoint, nil)
	req.SetBasicAuth(os.Getenv("WC_KEY"), os.Getenv("WC_SECRET"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during deleting product: %w", err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error during reading delete response body: %w", err)
	}
	return string(response), nil
}

func UpdateProduct(p Product) (string, error){
	err := validateProduct(p)
	if err != nil {
		return "", fmt.Errorf("Error during validating product: %w", err)
	}
	up := mapProductToUpProduct(p)
	wc_up_prdct_endpoint := os.Getenv("WC_URL") + "/wp-json/wc/v3/products/"
	wc_up_prdct_endpoint += fmt.Sprintf("%d", up.ID)
	wc_up_prdct_endpoint += "?consumer_key=" + os.Getenv("WP_KEY") + "&consumer_secret=" + os.Getenv("WP_SECRET")
	req, err := http.NewRequest("PATCH", wc_up_prdct_endpoint, bytes.NewBuffer([]byte(up.ToJson())))
	req.SetBasicAuth(os.Getenv("WC_KEY"), os.Getenv("WC_SECRET"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during posting product: %w", err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error during reading response body: %w", err)
	}
	return string(response), nil
}
func SendProduct(p Product) (string, error) {
	err := validateProduct(p)
	if err != nil {
		return "", fmt.Errorf("Error during validating product: %w", err)
	}
	wc_prdct_endpoint := os.Getenv("WC_URL") + "/wp-json/wc/v3/products/"
	wc_prdct_endpoint += "?consumer_key=" + os.Getenv("WP_KEY") + "&consumer_secret=" + os.Getenv("WP_SECRET")
	req, err := http.NewRequest("POST", wc_prdct_endpoint, bytes.NewBuffer([]byte(p.ToJson())))
	req.SetBasicAuth(os.Getenv("WC_KEY"), os.Getenv("WC_SECRET"))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error during posting product: %w", err)
	}
	defer resp.Body.Close()
	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error during reading response body: %w", err)
	}
	return string(response), nil
}
