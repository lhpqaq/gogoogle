package crawler

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

const defaultAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

type Google struct {
	url string
}

func (g Google) GetRaw(url string) {
	var rErr error
	ctx := context.Background()

	c := colly.NewCollector(colly.MaxDepth(1))
	c.UserAgent = defaultAgent
	_ = c.SetProxy("socks://127.0.0.1:7890")
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	nextPageLink := ""
	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
		if nextPageLink != "" {
			req, err := r.New("GET", nextPageLink, nil)
			if err == nil {
				q.AddRequest(req)
			}
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	rank := 1
	filteredRank := 1
	results := []Result{}
	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		// c.OnHTML("div[class='N54PNb BToiNc cvP2Ce']", func(e *colly.HTMLElement) {
		sel := e.DOM
		linkHref, _ := sel.Find("a").Attr("href")
		// fmt.Println(linkHref)
		linkText := strings.TrimSpace(linkHref)
		fmt.Println(linkText)
		// fmt.Println(sel.Html())
		// fmt.Println(sel.Find("div > div > div > span > a > h3").Text())
		titleText := strings.TrimSpace(sel.Find("div > div > div > span > a > h3").Text())
		descText := strings.TrimSpace(sel.Find("div > div > div > div:first-child > span:first-child").Text())
		fmt.Println(sel.Find("div > div > div > div:first-child > span:first-child").Text())
		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
			filteredRank += 1
		}
		fmt.Println("-----------------------------------------")
		// check if there is a next button at the end.
		// Added this selector as the Id is the same for every language checked on google.com .pt and .es the text changes but the id remains the same
		// nextPageHref, _ := sel.Find("a #pnnext").Attr("href")
		// nextPageLink = strings.TrimSpace(nextPageHref)
		// fmt.Println("here:", nextPageLink)
	})
	url += "&num=2"
	q.AddURL(url)
	q.Run(c)
	fmt.Println(rErr)
	// fmt.Println(results)
	// c.OnResponse(func(r *colly.Response) {
	// 	fmt.Printf("Response %s: %d bytes\n", r.Request.URL, len(r.Body))
	// 	fmt.Println()
	// })
	// c.Visit(url)
}
