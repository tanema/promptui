package promptui

import (
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/tanema/promptui/frmt"
	"github.com/tanema/promptui/screenbuf"
	"github.com/tanema/promptui/term"
)

type progress struct {
	title   string
	current float64
	total   float64
	err     error
	done    bool
}

// ProgressGroup tracks a group of progress bars
type ProgressGroup struct {
	items []*progress
	wg    sync.WaitGroup
}

// Progress will create a single progress bar and wait for it to finish
func Progress(total float64, fn func(func(float64)) error) error {
	group := NewProgressGroup()
	group.Go("", total, fn)
	return group.Wait()
}

// NewProgressGroup will create a new progress bar group the will track multiple bars
func NewProgressGroup() *ProgressGroup {
	return &ProgressGroup{}
}

// Go will add another bar to the group
func (pg *ProgressGroup) Go(title string, max float64, fn func(func(float64)) error) {
	pg.wg.Add(1)
	s := &progress{title: title, total: max}
	pg.items = append(pg.items, s)
	go func() {
		defer pg.wg.Done()
		s.err = fn(s.tick)
		s.done = true
	}()
}

// Wait will pause until all of the progress bars are complete
func (pg *ProgressGroup) Wait() error {
	done := false
	sb := screenbuf.New(os.Stdout)
	go func() {
		for !done {
			pg.render(sb)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	pg.wg.Wait()
	done = true
	pg.render(sb)
	return nil
}

func (pg *ProgressGroup) render(sb *screenbuf.ScreenBuf) {
	width, _, err := term.Size()
	if err != nil {
		return
	}
	sb.Reset()
	for _, item := range pg.items {
		sb.Write(item.render(width))
	}
	sb.Flush()
}

func (p *progress) tick(inc float64) {
	p.current = math.Min(p.current+inc, p.total)
	p.done = p.current == p.total
}

func (p *progress) render(width int) []byte {
	percent := p.current / p.total
	barwidth := width - len(p.title) - 8
	done := percent * float64(barwidth)
	status := "iconQ"
	if p.done && p.err == nil {
		status = "iconGood"
	} else if p.done {
		status = "iconBad"
	}
	return frmt.Render(
		fmt.Sprintf("{{%s}} {{.Title}} {{.Done | cyan}}{{.Rest}} {{.Percent}}%%", status),
		struct {
			Percent           int
			Title, Done, Rest string
		}{
			Title:   p.title,
			Percent: int(percent * 100),
			Done:    strings.Repeat("█", int(done)),
			Rest:    strings.Repeat("█", int(math.Max(float64(barwidth)-done, 0))),
		},
	)
}
