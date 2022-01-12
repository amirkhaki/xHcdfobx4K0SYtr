package main

import (
	"fmt"
	"io"
	"bytes"
	"golang.org/x/net/html"
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
		var buf bytes.Buffer
		w := io.Writer(&buf)

		err := html.Render(w, elms[i])

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(buf.String())
	}
}
