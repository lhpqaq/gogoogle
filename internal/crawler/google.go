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
	g.results = make(chan r.Result, 10)
	g.start = 0
	g.url = url
}

func (g *Google) ParseHTML() {
	var rErr error
	ctx := context.Background()
	c := colly.NewCollector(colly.MaxDepth(1))
	c.UserAgent = defaultAgent
	_ = c.SetProxy("socks://127.0.0.1:7890")
	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})
	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	rank := 1
	filteredRank := g.start + 1
	results := []r.Result{}
	c.OnHTML("div.g", func(e *colly.HTMLElement) {
		sel := e.DOM
		linkHref, _ := sel.Find("a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find("div > div > div > span > a > h3").Text())
		descText := strings.TrimSpace(sel.Find("div > div > div:nth-child(2) > div > span:nth-child(2)").Text())
		domainText := strings.TrimSpace(sel.Find("div > div > span > a > div > div > div > div:nth-child(1) > span").Text())
		timeText := strings.TrimSpace(sel.Find(".LEwnzc.Sqrs4e").Text())
		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := r.Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
				WebTime:     timeText,
				WebDomain:   domainText,
			}
			results = append(results, result)
			g.results <- result
			filteredRank += 1
		}
	})
	g.url += "&num=10"
	// g.url += "&num=10&lr=lang_zh-CN"
	g.url += fmt.Sprintf("&start=%d", g.start)
	q.AddURL(g.url)
	q.Run(c)
	if rErr != nil {
		fmt.Println(rErr)
	}
	result := r.Result{
		Rank: -1,
	}
	g.results <- result
}

func (g *Google) GetResult(re *r.Results) {
	re.Url = g.url
	go g.ParseHTML()
	num := 0
	for {
		res := <-g.results
		if res.Rank < 0 {
			g.start += num
			re.Parsing <- num
			num = 0
			cmd := <-re.Cmd
			switch cmd {
			case 1:
				{
					go g.ParseHTML()
				}
			}
		} else {
			re.Res[res.Rank] = res
			num += 1
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
