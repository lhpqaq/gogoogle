package result

import "fmt"

func (r *Results) Init() {
	r.Res = make(map[int]Result)
	r.Cmd = make(chan int)
	r.Parsing = make(chan int)
}

func (r *Results) Print() {
	start := 1
	for {
		num := <-r.Parsing
		fmt.Println(r.Res)
		for i := range num {
			re, ok := r.Res[start+i]
			if ok {
				fmt.Println(re)
			}
			start += 1
		}
		r.Cmd <- 1
	}
}
