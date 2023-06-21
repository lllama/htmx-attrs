package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	hxAttributes := make(map[string]bool)

	err := filepath.Walk(os.Args[1], func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".html" {
			text, err := os.Open(path)
			if err != nil {
				fmt.Printf("error reading the path %q \n", path)
				return err
			}
			doc := html.NewTokenizer(text)
			for {
				tt := doc.Next()
				switch tt {
				case html.ErrorToken:
					return err
				case html.StartTagToken, html.SelfClosingTagToken:
					token := doc.Token()
					for _, v := range token.Attr {
						if strings.HasPrefix(v.Key, "hx") {
							hxAttributes[v.Key] = true
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Print("error walking the path %q: %v\n")
		return
	}
	for key := range hxAttributes {
		fmt.Println(key)
	}
}
