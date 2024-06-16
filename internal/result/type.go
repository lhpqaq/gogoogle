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
}

type Results struct {
	Chan chan Result
	Num  int
}
