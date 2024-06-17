package crawler

import (
	r "gogoogle/internal/result"
)

type Crawler interface {
	Init(string)
	ParseHTML()
	GetResult(*r.Results)
	Print()
}

func Test(c Crawler, url string) {
	c.ParseHTML()
}

func Run(c Crawler, url string) {
	c.Init(url)
	go c.Print()
	c.ParseHTML()
}
