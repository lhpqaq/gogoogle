package result

import (
	"fmt"
	s "gogoogle/internal/show"
	"os"
	"os/exec"
	"strconv"
)

var gShow = s.Show{}

func (r *Results) Init() {
	r.Res = make(map[int]Result)
	r.Cmd = make(chan int)
	r.Parsing = make(chan int)
}

func (r *Results) Print() {
	start := 1
	for {
		num := <-r.Parsing
		r.Show(start, num)
		start += num
		r.Command()
	}
}

func (r *Results) Command() {
	for {
		fmt.Print("~ ")
		var cmd string
		fmt.Scanln(&cmd)
		switch cmd {
		case "n", "next":
			{
				r.Cmd <- 1
				return
			}
		case "q", "quit":
			{
				os.Exit(0)
			}
		case "o", "open":
			{
				exec.Command("open", r.Url).Start()
				break
			}
		default:
			{
				rank, err := strconv.Atoi(cmd)
				if err != nil {
					fmt.Println(err)
				}
				re, ok := r.Res[rank]
				if ok {
					exec.Command("open", re.URL).Start()
				}
			}
		}
	}
}

func (r *Results) Show(start int, num int) {
	// fmt.Println(start, num)
	for i := range num {
		re, ok := r.Res[start+i]
		if ok {
			re.show()
		}
	}
}

func (r *Result) show() {
	if gShow.Time == nil {
		gShow.Init()
	}
	gShow.Title.Printf("%d. %s", r.Rank, r.Title)
	gShow.Domain.Printf("   %s\n", r.WebDomain)
	gShow.Time.Printf("   %s ", r.WebTime)
	gShow.Desc.Printf("%s\n", r.Description)
	fmt.Println()
}
