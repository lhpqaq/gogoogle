package show

import (
	"github.com/fatih/color"
)

type Show struct {
	Title  *color.Color
	Time   *color.Color
	Desc   *color.Color
	Domain *color.Color
}

func (s *Show) Init() {
	s.Title = color.New(color.FgBlue).Add(color.Underline).Add(color.Bold)
	s.Time = color.New(90)
	s.Desc = color.New(color.FgWhite)
	s.Domain = color.New(color.FgGreen)
}
