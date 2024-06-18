package result

type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`

	// Time
	WebTime string `json:"webtime"`

	// Domain
	WebDomain string `json:"webdomain"`
}

type Results struct {
	Res     map[int]Result
	Url     string
	Parsing chan int
	Cmd     chan int
	Num     int
}
