package result

import (
	"fmt"
	"os/exec"
	"strconv"
)

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
			fmt.Printf("%2d. %s\n", start+i, re.Title)
			fmt.Printf("    %s%s\n", re.WebTime, re.WebDomain)
			fmt.Printf("    %s\n", re.Description)
		}
	}
}
