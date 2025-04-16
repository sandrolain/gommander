package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sandrolain/gommander/pkg/fs"
	"github.com/sandrolain/gommander/pkg/model"
)

var p *tea.Program

var leftWatcherSub *fs.DirWatcherSubscription
var rightWatcherSub *fs.DirWatcherSubscription

func main() {
	m := model.InitialModel(func(path string, cb func()) error {
		var err error
		leftWatcherSub, err = fs.SubscribeWatcher(leftWatcherSub, path, func(eventPath string, err error) {
			if err != nil {
				return
			}
			cb()
			if p != nil {
				p.Send(tea.FocusMsg{})
			}
		})
		return err
	}, func(path string, cb func()) error {
		var err error
		rightWatcherSub, err = fs.SubscribeWatcher(rightWatcherSub, path, func(eventPath string, err error) {
			if err != nil {
				return
			}
			cb()
			if p != nil {
				p.Send(tea.FocusMsg{})
			}
		})
		return err
	})

	p = tea.NewProgram(m, tea.WithMouseAllMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v", err)
		os.Exit(1)
	}
}
