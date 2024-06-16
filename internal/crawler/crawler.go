package crawler

type Crawler interface {
	Init(string)
	GetResult()
	Print()
}

func Test(c Crawler, url string) {
	c.GetResult()
}

func Run(c Crawler, url string) {
	c.Init(url)
	go c.Print()
	c.GetResult()
}
