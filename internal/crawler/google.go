package crawler

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	r "gogoogle/internal/result"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

const defaultAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

type Google struct {
	url     string
	results chan r.Result
	start   int
	topen   chan string
}

func (g *Google) Init(url string) {
	// fmt.Println(g.url)
	g.results = make(chan r.Result, 10)
	g.start = 0
	g.url = url
	// fmt.Println(g.url)
}

func (g *Google) ParseHTML() {
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
	filteredRank := g.start + 1
	results := []r.Result{}
	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		// c.OnHTML("div[class='N54PNb BToiNc cvP2Ce']", func(e *colly.HTMLElement) {
		sel := e.DOM
		linkHref, _ := sel.Find("a").Attr("href")
		// fmt.Println(linkHref)
		linkText := strings.TrimSpace(linkHref)
		// fmt.Println(linkText)
		// fmt.Println(sel.Html())
		// fmt.Println(sel.Find("div > div > div > span > a > h3").Text())
		titleText := strings.TrimSpace(sel.Find("div > div > div > span > a > h3").Text())
		descText := strings.TrimSpace(sel.Find("div > div > div > div:first-child > span:first-child").Text())
		// fmt.Println(sel.Find("div > div > div > div:first-child > span:first-child").Text())
		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := r.Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
			g.results <- result
			filteredRank += 1
		}
		// fmt.Println("-----------------------------------------")
	})
	g.url += "&num=10"
	g.url += fmt.Sprintf("&start=%d", g.start)
	// fmt.Println(g.url)
	q.AddURL(g.url)
	q.Run(c)
	fmt.Println(rErr)
	result := r.Result{
		Rank: -1,
	}
	g.results <- result
}

func (g *Google) GetResult(re *r.Results) {
	go g.ParseHTML()
	num := 0
	for {
		res := <-g.results
		num += 1
		if res.Rank < 0 {
			g.start += num
			re.Parsing <- num
			// fmt.Println(re.Res)
			num = 0
			cmd := <-re.Cmd
			switch cmd {
			case 1:
				{
					fmt.Println("cmd 1")
					go g.ParseHTML()
				}
			}
		} else {
			re.Res[res.Rank] = res
		}
	}
}

func (g *Google) Print() {
	num := 0
	for {
		res := <-g.results
		num += 1
		fmt.Println(res.Title)
		if res.Rank < 0 {
			g.start += num
			num = 0
			go g.ParseHTML()
		}
	}
}

func (g *Google) OpenUrl() {
	for {
		openurl := <-g.topen
		exec.Command("open", openurl).Start()
	}
}
