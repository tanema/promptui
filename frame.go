package promptui

import (
	"fmt"
	"strings"

	"github.com/tanema/promptui/frmt"
	"github.com/tanema/promptui/term"
)

const (
	DIV = "┣"
)

type Frame struct {
	Failure string
	Success string
	Timing  bool
}

func InFrame(title string, fn func() error) error {
	width, _, err := term.Size()
	if err != nil {
		return err
	}
	fmt.Println(string(frmt.Render("{{ . | cyan }}", "┏"+title+strings.Repeat("━", width-len(title)-1))))
	err = fn()
	if err != nil {
		fmt.Println(string(frmt.Render("{{ . | red }}", "┣Unexpected Error"+strings.Repeat("━", width-17))))
		fmt.Println(string(frmt.Render("{{ . | red }}", "┃"+err.Error())))
		fmt.Println(string(frmt.Render("{{ . | red }}", "┗"+strings.Repeat("━", width-1))))
		return err
	}
	fmt.Println(string(frmt.Render("{{ . | cyan }}", "┃")))
	fmt.Println(string(frmt.Render("{{ . | cyan }}", "┗"+strings.Repeat("━", width-1))))
	return nil
}
