package main

import (
	"fmt"
	r "gogoogle/internal/result"

	"gogoogle/internal/crawler"
	"gogoogle/internal/geturl"
	"os"
)

func main() {
	args := os.Args
	f := geturl.ParseArgs(args)
	fmt.Println("hello google")
	g := crawler.Google{}
	web, content, err := geturl.GetWeb(f)
	if f.Debug {
		fmt.Println("Web:", web)
	}
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	url := geturl.GetURL(web, geturl.ArrayToString(content, web.Delim))
	fmt.Println("URL:", url)
	res := r.Results{}
	res.Init()
	g.Init(url)
	go g.GetResult(&res)
	res.Print()
	// crawler.Run(&g, url)
}
