package main

import (
	"errors"
	"time"

	"github.com/tanema/promptui"
	//"github.com/tanema/promptui/input"
)

func main() {
	// (&input.Select{
	//	Label: "What's your text editor",
	//	Items: []string{"Vim", "Emacs", "Sublime", "VSCode", "Atom"},
	// }).Run()

	// (&input.Prompt{
	//	Label: "Whats your name",
	// }).Run()

	pgroup := promptui.NewProgressGroup()
	pgroup.Go("Git Clone", 100, func(tick func(val float64)) error {
		for i := 1; i <= 100; i++ {
			tick(1)
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	})
	pgroup.Go("Docker Image", 200, func(tick func(val float64)) error {
		for i := 1; i <= 100; i++ {
			tick(2)
			time.Sleep(10 * time.Millisecond)
			if i == 50 {
				return errors.New("connection error")
			}
		}
		return nil
	})
	pgroup.Go("Railgun Image", 50, func(tick func(val float64)) error {
		for i := 1; i <= 100; i++ {
			tick(1)
			time.Sleep(10 * time.Millisecond)
		}
		return nil
	})
	pgroup.Wait()
	// promptui.InFrame("test frame", func() error {
	//	return errors.New("foo")
	// })

	sgroup := promptui.NewSpinGroup()
	sgroup.Go("redis", func() error {
		time.Sleep(time.Second)
		return nil
	})
	sgroup.Go("mysql", func() error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})
	sgroup.Go("elasticsearch", func() error {
		time.Sleep(2 * time.Second)
		return errors.New("elasticseach failed to start")
	})
	sgroup.Wait()
}
