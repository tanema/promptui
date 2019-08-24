package promptui

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tanema/promptui/frmt"
	"github.com/tanema/promptui/screenbuf"
)

var glyphs = []rune("⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏")

type spinner struct {
	title string
	err   error
	done  bool
}

// SpinGroup keeps a group of spinners and their statuses
type SpinGroup struct {
	items   []*spinner
	current int
	on      bool
	wg      sync.WaitGroup
}

// Spinner creates a single spinner and waits for it to finish
func Spinner(title string, fn func() error) error {
	group := NewSpinGroup()
	group.Go(title, fn)
	return group.Wait()
}

// NewSpinGroup creates a new group of spinners to track multiple statuses
func NewSpinGroup() *SpinGroup {
	return &SpinGroup{}
}

// Go adds another process to the spin group
func (sg *SpinGroup) Go(title string, fn func() error) {
	sg.wg.Add(1)
	s := &spinner{title: title}
	sg.items = append(sg.items, s)
	go func() {
		defer sg.wg.Done()
		s.err = fn()
		s.done = true
	}()
}

// Wait will pause until all spinners are complete
func (sg *SpinGroup) Wait() error {
	done := false
	sb := screenbuf.New(os.Stdout)
	go func() {
		for !done {
			sg.next()
			sg.render(sb)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	sg.wg.Wait()
	done = true
	sg.render(sb)
	return nil
}

func (sg *SpinGroup) next() {
	sg.current++
	if sg.current >= len(glyphs) {
		sg.on = !sg.on
		sg.current = 0
	}
}

func (sg *SpinGroup) render(sb *screenbuf.ScreenBuf) {
	spn := frmt.Render("{{ . }} ", string(glyphs[sg.current]))
	if sg.on {
		spn = frmt.Render("{{ . | cyan }} ", string(glyphs[sg.current]))
	}
	sb.Reset()
	for _, item := range sg.items {
		if item.done && item.err == nil {
			sb.Write(frmt.Render("{{ iconGood }} {{ . }} ", item.title))
		} else if item.done {
			sb.Write(frmt.Render("{{ iconBad }} {{ . | red }} ", fmt.Sprintf("%s %v", item.title, item.err)))
		} else {
			sb.Write(append(spn, []byte(item.title)...))
		}
	}
	sb.Flush()
}
