package crawler

import "fmt"

type Crawler interface {
	GetRaw(string)
	Test()
}

type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`
}

func Test(c Crawler, url string) {
	c.GetRaw(url)
}

func (c Google) Test() {
	fmt.Println()
}
